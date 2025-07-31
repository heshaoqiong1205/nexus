package models

import (
	"time"

	"gorm.io/datatypes"
)

type Device struct {
	ID         string `gorm:"primaryKey"`
	SecretKey  string
	LicenseID  string
	Name       string
	ProductID  string
	GroupID    string
	Features   datatypes.JSON
	State      datatypes.JSON
	Version    string
	SDKVersion string
	IP         string
	Online     bool
	Location   *Point `gorm:"type:point"`
	CreatedAt  time.Time
	ActiveAt   time.Time
	UpdatedAt  time.Time `gorm:"autoUpdateTime:false"`
	Status     bool
}

type IDeviceModels interface {
	Get(id string) (Device, error)
	List(groupList []string, page int, pageSize int, orderBY *string) ([]Device, error)
	Create(device *Device) error
	Update(device *Device) error
	Delete(id string) error
}

type DeviceModels struct {
}

func (models *DeviceModels) Get(id string) (Device, error) {
	var device Device
	result := db.First(&device, "id = ?", id)
	if result.Error != nil {
		return Device{}, result.Error
	}
	return device, nil
}

func (models *DeviceModels) Count(groupList []string) (int64, error) {
	var total int64

	result := db.Model(&Device{}).Where("status = true and group_id IN ?", groupList).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *DeviceModels) List(groupList []string, page int, pageSize int, orderBY *string) ([]Device, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	var devices []Device

	if orderBY == nil {
		defaultOrder := "created_at desc"
		orderBY = &defaultOrder
	}
	result := db.Where("status = true and group_id IN ?", groupList).Order(*orderBY).Limit(limit).Offset(offset).Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (models *DeviceModels) Create(device *Device) error {
	result := db.Create(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *DeviceModels) Update(device *Device) error {
	result := db.Save(device)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *DeviceModels) Delete(id string) error {
	var device Device
	result := db.Delete(&device, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
