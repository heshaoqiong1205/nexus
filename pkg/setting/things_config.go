package setting

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type MQTT struct {
	Enabled            bool   `json:"enabled" ini:"enabled"`
	Url                string `json:"url" ini:"url"`
	ClientID           string `json:"client_id" ini:"client_id"`
	NodeID             string `json:"node_id" ini:"node_id"`
	Username           string `json:"username" ini:"username"`
	Password           string `json:"password" ini:"password"`
	Keepalive          int    `json:"keep_alive" ini:"keepalive"`
	CleanSession       bool   `json:"clean_session" ini:"clean_session"`
	ConnectTimeout     int    `json:"connect_timeout" ini:"connect_timeout"`
	QoS                int    `json:"qos" ini:"qos"`
	RPCTimeout         int    `json:"rpc_timeout" ini:"rpc_timeout"`
	SharedGroup        string `json:"shared_group" ini:"shared_group"`
	StatusTopic        string `json:"status_topic" ini:"status_topic"`
	StateReportedTopic string `json:"state_reported_topic" ini:"state_reported_topic"`
	StateDesiredTopic  string `json:"state_desired_topic" ini:"state_desired_topic"`
	MetricsTopic       string `json:"metrics_topic" ini:"metrics_topic"`
	EventTopic         string `json:"event_topic" ini:"event_topic"`
	RPCRequestTopic    string `json:"rpc_request_topic" ini:"rpc_request_topic"`
	RPCReplyTopic      string `json:"rpc_reply_topic" ini:"rpc_reply_topic"`
}

type EventConfig struct {
	Interval    int    `json:"interval" ini:"interval"`
	Medium      string `json:"medium" ini:"medium"`
	MaxDuration int    `json:"max_duration" ini:"max_duration"`
}

type RTCConfig struct {
	MaxDuration int `json:"max_duration" ini:"max_duration"`
	MaxRate     int `json:"max_rate" ini:"max_rate"`
}

type ThingsConfig struct {
	Mqtt  MQTT        `json:"mqtt"`
	Event EventConfig `json:"event"`
	Rtc   RTCConfig   `json:"rtc"`
}

var thingsConfig ThingsConfig

func init() {
	thingsConfig = initThingsConfig()
	log.Printf("Things configuration loaded: %+v", thingsConfig)
}

func GetThingsConfig() *ThingsConfig {
	return &thingsConfig
}

func initThingsConfig() ThingsConfig {
	configPaths := []string{
		"conf/things.ini",
		"../conf/things.ini",
		"../../conf/things.ini",
	}

	var cfg *ini.File
	var err error

	for _, path := range configPaths {
		if _, err = os.Stat(path); err == nil {
			cfg, err = ini.Load(path)
			if err == nil {
				log.Printf("Successfully loaded config file: %s", path)
				break
			}
		}
	}

	if err != nil {
		log.Printf("Warning: Could not load things.ini: %v. Using default configuration.", err)
		return getDefaultThingsConfig()
	}

	config := getDefaultThingsConfig()

	mqttSection := cfg.Section("mqtt")
	config.Mqtt.Enabled = mqttSection.Key("enabled").MustBool(config.Mqtt.Enabled)
	config.Mqtt.Url = stringOrDefault(mqttSection.Key("url").String(), config.Mqtt.Url)
	config.Mqtt.ClientID = mqttSection.Key("client_id").String()
	config.Mqtt.NodeID = mqttSection.Key("node_id").String()
	config.Mqtt.Username = mqttSection.Key("username").String()
	config.Mqtt.Password = mqttSection.Key("password").String()
	config.Mqtt.Keepalive = mqttSection.Key("keepalive").MustInt(config.Mqtt.Keepalive)
	config.Mqtt.CleanSession = mqttSection.Key("clean_session").MustBool(config.Mqtt.CleanSession)
	config.Mqtt.ConnectTimeout = mqttSection.Key("connect_timeout").MustInt(config.Mqtt.ConnectTimeout)
	config.Mqtt.QoS = mqttSection.Key("qos").MustInt(config.Mqtt.QoS)
	config.Mqtt.RPCTimeout = mqttSection.Key("rpc_timeout").MustInt(config.Mqtt.RPCTimeout)
	config.Mqtt.SharedGroup = mqttSection.Key("shared_group").String()
	config.Mqtt.StatusTopic = stringOrDefault(mqttSection.Key("status_topic").String(), config.Mqtt.StatusTopic)

	// Keep backward compatibility with the original key names while moving to the
	// new reported/desired/rpc topic split.
	config.Mqtt.StateReportedTopic = stringOrDefault(
		mqttSection.Key("state_reported_topic").String(),
		stringOrDefault(mqttSection.Key("state_topic").String(), config.Mqtt.StateReportedTopic),
	)
	config.Mqtt.StateDesiredTopic = stringOrDefault(mqttSection.Key("state_desired_topic").String(), config.Mqtt.StateDesiredTopic)
	config.Mqtt.MetricsTopic = stringOrDefault(mqttSection.Key("metrics_topic").String(), config.Mqtt.MetricsTopic)
	config.Mqtt.EventTopic = stringOrDefault(mqttSection.Key("event_topic").String(), config.Mqtt.EventTopic)
	config.Mqtt.RPCRequestTopic = stringOrDefault(
		mqttSection.Key("rpc_request_topic").String(),
		stringOrDefault(mqttSection.Key("service_topic").String(), config.Mqtt.RPCRequestTopic),
	)
	config.Mqtt.RPCReplyTopic = stringOrDefault(mqttSection.Key("rpc_reply_topic").String(), config.Mqtt.RPCReplyTopic)

	eventSection := cfg.Section("event")
	config.Event.Interval = eventSection.Key("interval").MustInt(config.Event.Interval)
	config.Event.Medium = stringOrDefault(eventSection.Key("medium").String(), config.Event.Medium)
	config.Event.MaxDuration = eventSection.Key("max_duration").MustInt(config.Event.MaxDuration)

	rtcSection := cfg.Section("rtc")
	config.Rtc.MaxDuration = rtcSection.Key("max_duration").MustInt(config.Rtc.MaxDuration)
	config.Rtc.MaxRate = rtcSection.Key("max_rate").MustInt(config.Rtc.MaxRate)

	return config
}

func getDefaultThingsConfig() ThingsConfig {
	return ThingsConfig{
		Mqtt: MQTT{
			Enabled:            true,
			Url:                "tcp://localhost:1883",
			Keepalive:          60,
			CleanSession:       true,
			ConnectTimeout:     5,
			QoS:                1,
			RPCTimeout:         10,
			StatusTopic:        "devices/{deviceId}/status",
			StateReportedTopic: "devices/{deviceId}/state/reported",
			StateDesiredTopic:  "devices/{deviceId}/state/desired",
			MetricsTopic:       "devices/{deviceId}/metrics",
			EventTopic:         "devices/{deviceId}/event",
			RPCRequestTopic:    "devices/{deviceId}/rpc/request",
			RPCReplyTopic:      "nexus/{nodeId}/rpc/reply",
		},
		Event: EventConfig{
			Interval:    60,
			Medium:      "mqtt",
			MaxDuration: 3600,
		},
		Rtc: RTCConfig{
			MaxDuration: 3600,
			MaxRate:     1024,
		},
	}
}

func stringOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
