package models

import "time"

type CloudRecordingPlan struct {
	ID        string `gorm:"primaryKey"`
	DeviceID  string
	Channel   int
	Enabled   bool
	Mode      string
	MFD       int
	Interval  int
	StorageID string
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
	Status    bool
}

type ICloudRecordingPlanModels interface {
	Get(id string) (CloudRecordingPlan, error)
	GetByDevice(deviceID string) ([]CloudRecordingPlan, error)
	Count() (int64, error)
	List(limt int, offset int, orderBY string) ([]CloudRecordingPlan, error)
	Create(plan *CloudRecordingPlan) error
	Update(plan *CloudRecordingPlan) error
	Delete(id string) error
}

type CloudRecordingPlanModels struct {
}

func (models *CloudRecordingPlanModels) Get(id string) (CloudRecordingPlan, error) {
	var plan CloudRecordingPlan
	result := db.Find(&plan, "id = ? and status = true", id)
	if result.Error != nil {
		return CloudRecordingPlan{}, result.Error
	}
	return plan, nil
}

func (models *CloudRecordingPlanModels) GetByDevice(deviceID string) ([]CloudRecordingPlan, error) {
	var plans []CloudRecordingPlan

	result := db.Find(&plans, "device_id = ? and status = true", deviceID)
	if result.Error != nil {
		return nil, result.Error
	}
	return plans, nil
}

func (models *CloudRecordingPlanModels) Count() (int64, error) {
	var total int64

	result := db.Model(&CloudRecordingPlan{}).Where("status = true").Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *CloudRecordingPlanModels) List(limt int, offset int, orderBY string) ([]CloudRecordingPlan, error) {
	var plans []CloudRecordingPlan

	result := db.Limit(limt).Offset(offset).Order(orderBY).Find(&plans, "status = true")
	if result.Error != nil {
		return nil, result.Error
	}
	return plans, nil
}

func (models *CloudRecordingPlanModels) Create(plan *CloudRecordingPlan) error {
	result := db.Create(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudRecordingPlanModels) Update(plan *CloudRecordingPlan) error {
	result := db.Save(plan)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudRecordingPlanModels) Delete(id string) error {
	result := db.Delete(&CloudRecordingPlan{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
