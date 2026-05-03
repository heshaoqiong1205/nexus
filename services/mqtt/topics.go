package mqtt

import (
	"fmt"
	"strings"
)

const (
	deviceIDPlaceholder = "{deviceId}"
	nodeIDPlaceholder   = "{nodeId}"
)

func RenderTopic(template string, values map[string]string) (string, error) {
	rendered := template
	for placeholder, value := range values {
		rendered = strings.ReplaceAll(rendered, placeholder, value)
	}

	if strings.Contains(rendered, "{") || strings.Contains(rendered, "}") {
		return "", fmt.Errorf("unresolved topic template: %s", template)
	}
	return rendered, nil
}

func RenderDeviceTopic(template string, deviceID string, nodeID string) (string, error) {
	values := map[string]string{
		deviceIDPlaceholder: deviceID,
		nodeIDPlaceholder:   nodeID,
	}
	return RenderTopic(template, values)
}

func SubscriptionTopic(template string, nodeID string, wildcardDevice bool, sharedGroup string) (string, error) {
	deviceValue := ""
	if wildcardDevice {
		deviceValue = "+"
	}

	rendered, err := RenderTopic(template, map[string]string{
		deviceIDPlaceholder: deviceValue,
		nodeIDPlaceholder:   nodeID,
	})
	if err != nil {
		return "", err
	}

	if sharedGroup != "" && wildcardDevice {
		return fmt.Sprintf("$share/%s/%s", sharedGroup, rendered), nil
	}
	return rendered, nil
}

func MatchTopic(template string, topic string, nodeID string) (map[string]string, bool) {
	templateParts := splitTopic(template)
	topicParts := splitTopic(topic)
	if len(templateParts) != len(topicParts) {
		return nil, false
	}

	params := make(map[string]string)
	for idx, part := range templateParts {
		switch part {
		case deviceIDPlaceholder:
			params["deviceId"] = topicParts[idx]
		case nodeIDPlaceholder:
			if nodeID != "" && topicParts[idx] != nodeID {
				return nil, false
			}
			params["nodeId"] = topicParts[idx]
		default:
			if part != topicParts[idx] {
				return nil, false
			}
		}
	}
	return params, true
}

func splitTopic(topic string) []string {
	return strings.Split(strings.Trim(topic, "/"), "/")
}
