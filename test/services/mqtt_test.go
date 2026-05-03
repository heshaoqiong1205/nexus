package services_testing

import (
	"testing"

	mqtt "nexus/services/mqtt"
)

func TestRenderDeviceTopic(t *testing.T) {
	topic, err := mqtt.RenderDeviceTopic("devices/{deviceId}/state/reported", "dev-001", "node-a")
	if err != nil {
		t.Fatalf("RenderDeviceTopic returned error: %v", err)
	}
	if topic != "devices/dev-001/state/reported" {
		t.Fatalf("unexpected topic: %s", topic)
	}
}

func TestSubscriptionTopicWithSharedGroup(t *testing.T) {
	topic, err := mqtt.SubscriptionTopic("devices/{deviceId}/metrics", "node-a", true, "nexus")
	if err != nil {
		t.Fatalf("SubscriptionTopic returned error: %v", err)
	}
	expected := "$share/nexus/devices/+/metrics"
	if topic != expected {
		t.Fatalf("expected %s, got %s", expected, topic)
	}
}

func TestMatchTopic(t *testing.T) {
	params, ok := mqtt.MatchTopic("devices/{deviceId}/event", "devices/dev-002/event", "node-a")
	if !ok {
		t.Fatal("expected topic to match")
	}
	if params["deviceId"] != "dev-002" {
		t.Fatalf("expected deviceId dev-002, got %s", params["deviceId"])
	}
}

func TestMatchTopicWithNodeID(t *testing.T) {
	params, ok := mqtt.MatchTopic("nexus/{nodeId}/rpc/reply", "nexus/node-a/rpc/reply", "node-a")
	if !ok {
		t.Fatal("expected topic to match")
	}
	if params["nodeId"] != "node-a" {
		t.Fatalf("expected nodeId node-a, got %s", params["nodeId"])
	}
}
