package mqtt

import (
	"log"
	"time"

	"nexus/services/things"
)

type Dispatcher interface {
	HandleStatus(deviceID string, status StatusMessage) error
	HandleState(deviceID string, state things.State) error
	HandleEvent(deviceID string, event things.Event) error
	HandleMetrics(deviceID string, payload []byte) error
}

type ServiceDispatcher struct {
	thingService *things.ThingService
}

func NewServiceDispatcher(thingService *things.ThingService) *ServiceDispatcher {
	return &ServiceDispatcher{
		thingService: thingService,
	}
}

func (d *ServiceDispatcher) HandleStatus(deviceID string, status StatusMessage) error {
	updatedAt := time.Now()
	if status.Timestamp > 0 {
		updatedAt = time.UnixMilli(status.Timestamp)
	}
	return d.thingService.HandleStatus(deviceID, status.Online, updatedAt)
}

func (d *ServiceDispatcher) HandleState(deviceID string, state things.State) error {
	return d.thingService.HandleState(deviceID, state)
}

func (d *ServiceDispatcher) HandleEvent(deviceID string, event things.Event) error {
	return d.thingService.HandleEvent(deviceID, event)
}

func (d *ServiceDispatcher) HandleMetrics(deviceID string, payload []byte) error {
	log.Printf("Received metrics from device %s: %s", deviceID, string(payload))
	return nil
}
