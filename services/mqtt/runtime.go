package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"nexus/pkg/setting"
	"nexus/services/things"
)

type rpcResult struct {
	response *things.ServiceResponse
	err      error
}

type Runtime struct {
	cfg        *setting.ThingsConfig
	dispatcher Dispatcher
	cm         *autopaho.ConnectionManager
	cancel     context.CancelFunc
	nodeID     string

	mu      sync.Mutex
	pending map[string]chan rpcResult
}

func NewRuntime(cfg *setting.ThingsConfig, dispatcher Dispatcher) (*Runtime, error) {
	if cfg == nil {
		return nil, errors.New("things config cannot be nil")
	}
	if dispatcher == nil {
		return nil, errors.New("dispatcher cannot be nil")
	}

	brokerURL, err := normalizeBrokerURL(cfg.Mqtt.Url)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(brokerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid mqtt broker url %q: %w", brokerURL, err)
	}

	nodeID := cfg.Mqtt.NodeID
	if strings.TrimSpace(nodeID) == "" {
		hostname, hostErr := os.Hostname()
		if hostErr != nil || hostname == "" {
			hostname = "localhost"
		}
		nodeID = fmt.Sprintf("%s-%d", hostname, os.Getpid())
	}

	rt := &Runtime{
		cfg:        cfg,
		dispatcher: dispatcher,
		nodeID:     nodeID,
		pending:    make(map[string]chan rpcResult),
	}

	ctx, cancel := context.WithCancel(context.Background())
	rt.cancel = cancel

	clientConfig := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{parsedURL},
		KeepAlive:                     uint16(cfg.Mqtt.Keepalive),
		CleanStartOnInitialConnection: cfg.Mqtt.CleanSession,
		SessionExpiryInterval:         sessionExpiryInterval(cfg.Mqtt.CleanSession),
		ReconnectBackoff:              autopaho.NewConstantBackoff(2 * time.Second),
		ConnectTimeout:                time.Duration(cfg.Mqtt.ConnectTimeout) * time.Second,
		OnConnectError: func(err error) {
			log.Printf("MQTT connect error: %v", err)
		},
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			rt.onConnectionUp(cm)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: defaultClientID(cfg.Mqtt.ClientID, nodeID),
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				rt.handlePublishReceived,
			},
			OnClientError: func(err error) {
				log.Printf("MQTT client error: %v", err)
			},
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d != nil && d.Properties != nil && d.Properties.ReasonString != "" {
					log.Printf("MQTT server disconnect: %s", d.Properties.ReasonString)
					return
				}
				if d != nil {
					log.Printf("MQTT server disconnect with reason code %d", d.ReasonCode)
				}
			},
		},
	}

	if cfg.Mqtt.Username != "" {
		clientConfig.ConnectUsername = cfg.Mqtt.Username
	}
	if cfg.Mqtt.Password != "" {
		clientConfig.ConnectPassword = []byte(cfg.Mqtt.Password)
	}

	cm, err := autopaho.NewConnection(ctx, clientConfig)
	if err != nil {
		cancel()
		return nil, err
	}
	rt.cm = cm
	return rt, nil
}

func (r *Runtime) Start() error {
	if r.cm == nil {
		return errors.New("mqtt connection manager is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Mqtt.ConnectTimeout)*time.Second)
	defer cancel()
	return r.cm.AwaitConnection(ctx)
}

func (r *Runtime) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
	if r.cm != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.cm.Disconnect(ctx); err != nil {
			log.Printf("MQTT disconnect error: %v", err)
		}
	}
}

func (r *Runtime) Call(request *things.ServiceRequest) (*things.ServiceResponse, error) {
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}
	if request.DeviceID == "" {
		return nil, errors.New("device id cannot be empty")
	}
	if request.TransactionID == "" {
		return nil, errors.New("transaction id cannot be empty")
	}
	if r.cm == nil {
		return nil, errors.New("mqtt connection manager is not available")
	}

	requestTopic, err := RenderDeviceTopic(r.cfg.Mqtt.RPCRequestTopic, request.DeviceID, r.nodeID)
	if err != nil {
		return nil, err
	}
	replyTopic, err := RenderDeviceTopic(r.cfg.Mqtt.RPCReplyTopic, request.DeviceID, r.nodeID)
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(r.cfg.Mqtt.RPCTimeout) * time.Second
	payload, err := json.Marshal(RPCRequestEnvelope{
		ID:        request.TransactionID,
		DeviceID:  request.DeviceID,
		Method:    request.Method,
		Data:      json.RawMessage(request.Data),
		Timestamp: time.Now().UnixMilli(),
		TimeoutMS: timeout.Milliseconds(),
	})
	if err != nil {
		return nil, err
	}

	resultCh := make(chan rpcResult, 1)
	if err := r.registerPending(request.TransactionID, resultCh); err != nil {
		return nil, err
	}
	defer r.unregisterPending(request.TransactionID)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	messageExpiry := uint32(timeout / time.Second)
	if messageExpiry == 0 {
		messageExpiry = 1
	}

	if _, err := r.cm.Publish(ctx, &paho.Publish{
		QoS:     byte(r.cfg.Mqtt.QoS),
		Topic:   requestTopic,
		Payload: payload,
		Properties: &paho.PublishProperties{
			ResponseTopic:   replyTopic,
			CorrelationData: []byte(request.TransactionID),
			MessageExpiry:   &messageExpiry,
		},
	}); err != nil {
		return nil, err
	}

	select {
	case result := <-resultCh:
		return result.response, result.err
	case <-ctx.Done():
		return nil, fmt.Errorf("rpc timeout waiting for device %s", request.DeviceID)
	}
}

func (r *Runtime) onConnectionUp(cm *autopaho.ConnectionManager) {
	log.Printf("MQTT connected to broker %s as node %s", r.cfg.Mqtt.Url, r.nodeID)

	subscriptions := []struct {
		template       string
		wildcardDevice bool
		shared         bool
	}{
		{template: r.cfg.Mqtt.StatusTopic, wildcardDevice: true, shared: true},
		{template: r.cfg.Mqtt.StateReportedTopic, wildcardDevice: true, shared: true},
		{template: r.cfg.Mqtt.EventTopic, wildcardDevice: true, shared: true},
		{template: r.cfg.Mqtt.MetricsTopic, wildcardDevice: true, shared: true},
		{template: r.cfg.Mqtt.RPCReplyTopic, wildcardDevice: false, shared: false},
	}

	for _, sub := range subscriptions {
		topic, err := SubscriptionTopic(sub.template, r.nodeID, sub.wildcardDevice, sharedGroupName(r.cfg.Mqtt.SharedGroup, sub.shared))
		if err != nil {
			log.Printf("Failed to build subscription topic from %s: %v", sub.template, err)
			continue
		}
		if err := subscribe(cm, topic, byte(r.cfg.Mqtt.QoS)); err != nil {
			log.Printf("Failed to subscribe topic %s: %v", topic, err)
		}
	}
}

func (r *Runtime) handlePublishReceived(pr paho.PublishReceived) (bool, error) {
	if pr.Packet == nil {
		return false, nil
	}

	switch {
	case matchTemplate(r.cfg.Mqtt.RPCReplyTopic, pr.Packet.Topic, r.nodeID):
		r.handleRPCReplyMessage(pr.Packet)
		return true, nil
	case matchTemplate(r.cfg.Mqtt.StatusTopic, pr.Packet.Topic, r.nodeID):
		r.handleStatusMessage(pr.Packet)
		return true, nil
	case matchTemplate(r.cfg.Mqtt.StateReportedTopic, pr.Packet.Topic, r.nodeID):
		r.handleStateMessage(pr.Packet)
		return true, nil
	case matchTemplate(r.cfg.Mqtt.EventTopic, pr.Packet.Topic, r.nodeID):
		r.handleEventMessage(pr.Packet)
		return true, nil
	case matchTemplate(r.cfg.Mqtt.MetricsTopic, pr.Packet.Topic, r.nodeID):
		r.handleMetricsMessage(pr.Packet)
		return true, nil
	default:
		log.Printf("MQTT received unmatched message on topic %s", pr.Packet.Topic)
		return false, nil
	}
}

func (r *Runtime) handleStatusMessage(message *paho.Publish) {
	params, ok := MatchTopic(r.cfg.Mqtt.StatusTopic, message.Topic, r.nodeID)
	if !ok {
		return
	}

	var status StatusMessage
	if err := json.Unmarshal(message.Payload, &status); err != nil {
		log.Printf("Failed to decode device status from topic %s: %v", message.Topic, err)
		return
	}

	go func(deviceID string, current StatusMessage) {
		r.logDispatchError("status", deviceID, r.dispatcher.HandleStatus(deviceID, current))
	}(params["deviceId"], status)
}

func (r *Runtime) handleStateMessage(message *paho.Publish) {
	params, ok := MatchTopic(r.cfg.Mqtt.StateReportedTopic, message.Topic, r.nodeID)
	if !ok {
		return
	}

	var state things.State
	if err := json.Unmarshal(message.Payload, &state); err != nil {
		log.Printf("Failed to decode device state from topic %s: %v", message.Topic, err)
		return
	}

	go func(deviceID string, current things.State) {
		r.logDispatchError("state", deviceID, r.dispatcher.HandleState(deviceID, current))
	}(params["deviceId"], state)
}

func (r *Runtime) handleEventMessage(message *paho.Publish) {
	params, ok := MatchTopic(r.cfg.Mqtt.EventTopic, message.Topic, r.nodeID)
	if !ok {
		return
	}

	var event things.Event
	if err := json.Unmarshal(message.Payload, &event); err != nil {
		log.Printf("Failed to decode device event from topic %s: %v", message.Topic, err)
		return
	}

	go func(deviceID string, current things.Event) {
		r.logDispatchError("event", deviceID, r.dispatcher.HandleEvent(deviceID, current))
	}(params["deviceId"], event)
}

func (r *Runtime) handleMetricsMessage(message *paho.Publish) {
	params, ok := MatchTopic(r.cfg.Mqtt.MetricsTopic, message.Topic, r.nodeID)
	if !ok {
		return
	}

	payload := cloneBytes(message.Payload)
	go func(deviceID string, currentPayload []byte) {
		r.logDispatchError("metrics", deviceID, r.dispatcher.HandleMetrics(deviceID, currentPayload))
	}(params["deviceId"], payload)
}

func (r *Runtime) handleRPCReplyMessage(message *paho.Publish) {
	requestID := ""
	if message.Properties != nil && len(message.Properties.CorrelationData) > 0 {
		requestID = string(message.Properties.CorrelationData)
	}

	var response RPCResponseEnvelope
	if len(message.Payload) > 0 {
		if err := json.Unmarshal(message.Payload, &response); err != nil {
			log.Printf("Failed to decode rpc reply from topic %s: %v", message.Topic, err)
			return
		}
	}

	if requestID == "" {
		requestID = response.ID
	}
	if requestID == "" {
		log.Printf("Discarding rpc reply without correlation id on topic %s", message.Topic)
		return
	}

	result := rpcResult{
		response: &things.ServiceResponse{
			DeviceID:      response.DeviceID,
			TransactionID: requestID,
			Data:          cloneBytes(response.Data),
		},
	}
	if len(result.response.Data) == 0 {
		result.response.Data = cloneBytes(message.Payload)
	}
	if response.Error != "" {
		result.err = errors.New(response.Error)
	}

	if !r.resolvePending(requestID, result) {
		log.Printf("Received rpc reply for unknown request %s", requestID)
	}
}

func (r *Runtime) registerPending(requestID string, resultCh chan rpcResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pending[requestID]; exists {
		return fmt.Errorf("rpc request %s is already pending", requestID)
	}
	r.pending[requestID] = resultCh
	return nil
}

func (r *Runtime) unregisterPending(requestID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.pending, requestID)
}

func (r *Runtime) resolvePending(requestID string, result rpcResult) bool {
	r.mu.Lock()
	resultCh, exists := r.pending[requestID]
	if exists {
		delete(r.pending, requestID)
	}
	r.mu.Unlock()
	if !exists {
		return false
	}

	select {
	case resultCh <- result:
	default:
	}
	return true
}

func (r *Runtime) logDispatchError(kind string, deviceID string, err error) {
	if err != nil {
		log.Printf("Failed to handle %s message for device %s: %v", kind, deviceID, err)
	}
}

func subscribe(cm *autopaho.ConnectionManager, topic string, qos byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cm.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{Topic: topic, QoS: qos},
		},
	})
	return err
}

func normalizeBrokerURL(rawURL string) (string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", errors.New("mqtt broker url cannot be empty")
	}

	switch {
	case strings.HasPrefix(rawURL, "tcp://"):
		return "mqtt://" + strings.TrimPrefix(rawURL, "tcp://"), nil
	case strings.HasPrefix(rawURL, "ssl://"):
		return "tls://" + strings.TrimPrefix(rawURL, "ssl://"), nil
	case strings.HasPrefix(rawURL, "mqtts://"):
		return "tls://" + strings.TrimPrefix(rawURL, "mqtts://"), nil
	default:
		return rawURL, nil
	}
}

func defaultClientID(clientID string, nodeID string) string {
	if strings.TrimSpace(clientID) != "" {
		return clientID
	}
	return fmt.Sprintf("nexus-%s", nodeID)
}

func sharedGroupName(group string, enabled bool) string {
	if !enabled {
		return ""
	}
	return group
}

func sessionExpiryInterval(cleanSession bool) uint32 {
	if cleanSession {
		return 0
	}
	return 600
}

func cloneBytes(payload []byte) []byte {
	if payload == nil {
		return nil
	}
	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return cloned
}

func matchTemplate(template string, topic string, nodeID string) bool {
	_, ok := MatchTopic(template, topic, nodeID)
	return ok
}
