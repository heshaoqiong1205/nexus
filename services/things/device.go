package things

import (
	"encoding/json"
	"nexus/models"
)

type IoTDevice struct {
	ID         string   `json:"id" gorm:"primaryKey"`
	Name       string   `json:"name"`
	ProductID  string   `json:"product_id"`
	GroupID    string   `json:"group_id"`
	Features   Features `json:"features"`
	State      state    `json:"state"`
	Version    string   `json:"version"`
	SDKVersion string   `json:"sdk_version"`
	IP         string   `json:"ip"`
	Online     bool     `json:"online"`
	Location   location `json:"location"`
	CreatedAt  int64    `json:"created_at"`
	UpdatedAt  int64    `json:"updated_at"`
	Status     bool     `json:"status"`
}

func NewIoTDevice(device models.Device) (*IoTDevice, error) {
	var state state
	var features Features
	err := json.Unmarshal(device.Features, &features)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(device.State, &state)
	if err != nil {
		return nil, err
	}
	return &IoTDevice{
		ID:         device.ID,
		Name:       device.Name,
		ProductID:  device.ProductID,
		GroupID:    device.GroupID,
		Features:   features,
		State:      state,
		Version:    device.Version,
		SDKVersion: device.SDKVersion,
		IP:         device.IP,
		Online:     device.Online,
		Status:     device.Status,
	}, nil
}

func (d *IoTDevice) ToDevice() (*models.Device, error) {
	if err := d.State.Validate(); err != nil {
		return nil, err
	}
	if err := d.Location.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}
