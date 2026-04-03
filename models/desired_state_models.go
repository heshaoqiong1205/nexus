package models

import (
	"errors"
	"time"

	"gorm.io/datatypes"
)

type DesiredState struct {
	ID            string         `gorm:"primaryKey" json:"id"` // DeviceID as primary key
	State         datatypes.JSON `gorm:"column:state;type:text" json:"state"`
	Version       int            `gorm:"default:1" json:"version"`
	Status        bool           `gorm:"default:true" json:"status"`
	LastDesiredID int64          `json:"last_desired_id"`
	ConfirmedAt   *time.Time     `json:"confirmed_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type IDesiredStateModels interface {
	Create(desiredState *DesiredState) error
	Get(deviceID string) (DesiredState, error)
	Update(desiredState *DesiredState) error
	UpdateWithVersion(desiredState *DesiredState, version int) error
	Delete(deviceID string) error
	List(deviceIDs []string, page int, pageSize int, orderBy *string) ([]DesiredState, error)
	Count(deviceIDs []string) (int64, error)
}

type DesiredStateModels struct {
}

func (d *DesiredStateModels) Create(desiredState *DesiredState) error {
	if desiredState.ID == "" {
		return errors.New("desired state id cannot be empty")
	}
	if desiredState.State == nil {
		desiredState.State = datatypes.JSON("{}")
	}
	desiredState.Status = true
	return db.Create(desiredState).Error
}

func (d *DesiredStateModels) Get(deviceID string) (DesiredState, error) {
	var desiredState DesiredState
	err := db.Where("id = ?", deviceID).First(&desiredState).Error
	return desiredState, err
}

func (d *DesiredStateModels) Update(desiredState *DesiredState) error {
	return db.Save(desiredState).Error
}

func (d *DesiredStateModels) Delete(deviceID string) error {
	return db.Where("id = ?", deviceID).Delete(&DesiredState{}).Error
}

func (d *DesiredStateModels) UpdateWithVersion(desiredState *DesiredState, version int) error {
	result := db.Model(&DesiredState{}).Where("id = ? and version = ?", desiredState.ID, version).Updates(desiredState)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("desired state version conflict")
	}
	return nil
}

func (d *DesiredStateModels) List(deviceIDs []string, page int, pageSize int, orderBy *string) ([]DesiredState, error) {
	var desiredStates []DesiredState

	limit := pageSize
	offset := (page - 1) * pageSize
	query := db.Model(&DesiredState{})

	if len(deviceIDs) > 0 {
		query = query.Where("id IN ?", deviceIDs)
	}

	if orderBy == nil {
		defaultOrder := "created_at DESC"
		orderBy = &defaultOrder
	}
	result := query.Order(*orderBy).Limit(limit).Offset(offset).Find(&desiredStates)
	return desiredStates, result.Error
}

func (d *DesiredStateModels) Count(deviceIDs []string) (int64, error) {
	var count int64
	query := db.Model(&DesiredState{})

	if len(deviceIDs) > 0 {
		query = query.Where("id IN ?", deviceIDs)
	}
	result := query.Count(&count)
	return count, result.Error
}
