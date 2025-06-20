package models

import (
	"time"

	"gorm.io/datatypes"
)

type CloudRecording struct {
	ID               string `gorm:"primaryKey"`
	DeviceID         string
	Channel          int
	FragmentDuration int
	BucketID         string
	Prefix           string
	Fragments        datatypes.JSON
	State            string
	BeginTime        time.Time
	EndTime          time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time `gorm:"autoUpdateTime:false"`
	Status           bool
}

type ICloudRecordingModels interface {
	Get(id string) (CloudRecording, error)
	GetByRecord(deviceID string, channelIdx int, bucketID string, prefix string) (CloudRecording, error)
	Count(deviceID string, channelIdx int, begin, end time.Time) (int64, error)
	List(deviceID string, channelIdx int, begin, end time.Time, limt int, offset int) ([]CloudRecording, error)
	Create(recording *CloudRecording) error
	Update(recording *CloudRecording) error
	Delete(id string) error
}

type CloudRecordingModels struct {
}

func (models *CloudRecordingModels) Get(id string) (CloudRecording, error) {
	var recording CloudRecording
	result := db.Find(&recording, "id = ?", id)
	if result.Error != nil {
		return CloudRecording{}, result.Error
	}
	return recording, nil
}

func (models *CloudRecordingModels) GetByRecord(deviceID string, channel int, bucketID string, prefix string) (CloudRecording, error) {
	var recording CloudRecording

	result := db.Find(&recording, "device_id = ? and channel = ? and bucket_id = ? and prefix = ? and status = true",
		deviceID, channel, bucketID, prefix)
	if result.Error != nil {
		return CloudRecording{}, result.Error
	}
	return recording, nil
}

func (models *CloudRecordingModels) Count(deviceID string, channel int, begin, end time.Time) (int64, error) {
	var total int64

	result := db.Model(&CloudRecording{}).
		Where("device_id = ? and channel = ? and begin_time >= ? and begin_time =< ? and status = true", deviceID, channel, begin, end).
		Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *CloudRecordingModels) List(deviceID string, channel int, begin, end time.Time, limt int, offset int) ([]CloudRecording, error) {
	var recordings []CloudRecording

	result := db.Limit(limt).Offset(offset).Order("begin_time desc").Find(&recordings,
		"device_id = ? and channel = ? and begin_time >= ? and begin_time =< ? and status = true", deviceID, channel, begin, end)
	if result.Error != nil {
		return nil, result.Error
	}
	return recordings, nil
}

func (models *CloudRecordingModels) Create(recording *CloudRecording) error {
	result := db.Create(recording)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudRecordingModels) Update(recording *CloudRecording) error {
	result := db.Save(recording)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudRecordingModels) Delete(id string) error {
	var recording CloudRecording
	result := db.Delete(&recording, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
