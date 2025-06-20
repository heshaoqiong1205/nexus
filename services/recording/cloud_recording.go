package recording

import (
	"encoding/json"
	"errors"
	"nexus/models"
	"nexus/services/storage"
	"time"

	"log"

	"gorm.io/gorm"
)

const (
	RecordingModeManual   = "manual"
	RecordingModeSchedule = "schedule"
	RecordingModeMotion   = "motion"
)

type StorageInfo struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

type Fragment struct {
	Duartion int    `json:"duration"` // Duration in seconds
	Index    string `json:"index"`    // Unique identifier for the fragment
}

type BeginRecordingRequest struct {
	DeviceID  string      `json:"device_id"`
	Channel   int         `json:"channel"`
	Storage   StorageInfo `json:"storage"`
	StartTime time.Time   `json:"start_time"`
	MFD       int         `json:"mfd"`
}

type EndRecordingRequest struct {
	DeviceID  string      `json:"device_id"`
	Channel   int         `json:"channel"`
	Storage   StorageInfo `json:"storage"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	MFD       int         `json:"mfd"`
	Fragments []Fragment  `json:"fragments"`
}

type RecordingConfig struct {
	Enabled  bool   `json:"enabled"`
	Mode     string `json:"mode"`
	MFD      int    `json:"mfd"`
	Interval int    `json:"interval"`
}

type CloudRecordingService struct {
	cloudRecordingModels models.ICloudRecordingModels
	cloudStorageService  storage.ICloudStorageService
	bucketService        storage.IBucketService
}

func NewCloudRecordingService(cloudRecordingModels models.ICloudRecordingModels,
	cloudStorageService storage.ICloudStorageService, bucketService storage.IBucketService) *CloudRecordingService {
	return &CloudRecordingService{
		cloudRecordingModels,
		cloudStorageService,
		bucketService,
	}
}

func (s *CloudRecordingService) BeginRecording(req BeginRecordingRequest) error {
	// Validate the request
	if req.DeviceID == "" {
		return errors.New("invalid request: device ID  must be provided")
	}

	if req.StartTime.Before(time.Now().Add(5*time.Minute)) || req.StartTime.After(time.Now().Add(5*time.Minute)) {
		return errors.New("invalid request: start time must be within 5 minutes from now")
	}

	if req.MFD < 5 || req.MFD > 60 {
		return errors.New("invalid request: MFD must be between 5 and 60 seconds")
	}

	if req.Storage.Prefix == "" {
		return errors.New("invalid request: storage prefix must be provided")
	}
	bucket, err := s.bucketService.GetBucketByName(req.Storage.Bucket)
	if err != nil {
		return errors.New("invalid request: bucket not found")

	}
	if req.Channel < 0 {
		return errors.New("invalid request: channel index must be a non-negative integer")
	}

	// Here you would implement the logic to start recording
	// This could involve interacting with a cloud storage service, etc.
	// For now, we will just log the request
	log.Printf("Begin recording for device %s on channel %d at %s", req.DeviceID, req.Channel, req.StartTime)

	cloudRecording := models.CloudRecording{
		ID:               req.DeviceID,
		Channel:          req.Channel,
		FragmentDuration: req.MFD,
		BucketID:         bucket.ID,
		Prefix:           req.Storage.Prefix,
		Fragments:        nil, // This would be populated with actual fragment data
		State:            "recording",
		BeginTime:        req.StartTime,
		Status:           true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	s.cloudRecordingModels.Create(&cloudRecording)

	return nil
}

func (s *CloudRecordingService) EndRecording(req EndRecordingRequest) error {
	// Validate the request
	if req.DeviceID == "" {
		return errors.New("invalid request: device ID must be provided")
	}

	if req.StartTime.Before(time.Now().Add(5*time.Minute)) || req.StartTime.After(time.Now().Add(5*time.Minute)) {
		return errors.New("invalid request: start time must be within 5 minutes from now")
	}

	if req.EndTime.Before(req.StartTime) {
		return errors.New("invalid request: end time must be after start time")
	}

	if req.MFD < 5 || req.MFD > 60 {
		return errors.New("invalid request: MFD must be between 5 and 60 seconds")
	}

	if len(req.Fragments) <= 0 {
		return errors.New("invalid request: at least one fragment must be provided")
	}
	fragments, err := json.Marshal(req.Fragments)
	if err != nil {
		return errors.New("invalid request: fragments must be valid JSON")
	}

	if req.Channel < 0 {
		return errors.New("invalid request: channel index must be a non-negative integer")
	}

	bucket, err := s.bucketService.GetBucketByName(req.Storage.Bucket)
	if err != nil {
		return errors.New("invalid request: bucket not found")

	}
	log.Printf("End recording for device %s on channel %d from %s to %s", req.DeviceID, req.Channel, req.StartTime, req.EndTime)

	recording, err := s.cloudRecordingModels.GetByRecord(req.DeviceID, req.Channel, bucket.ID, req.Storage.Prefix)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cloudRecording := models.CloudRecording{
				ID:               req.DeviceID,
				Channel:          req.Channel,
				FragmentDuration: req.MFD,
				BucketID:         bucket.ID,
				Prefix:           req.Storage.Prefix,
				Fragments:        fragments,
				State:            "complete",
				BeginTime:        req.StartTime,
				EndTime:          req.EndTime,
				Status:           true,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			s.cloudRecordingModels.Create(&cloudRecording)
		}
		return err
	}

	recording.State = "complete"
	recording.EndTime = req.EndTime
	recording.UpdatedAt = time.Now()
	recording.Fragments = fragments
	s.cloudRecordingModels.Update(&recording)

	return nil
}

func (s *CloudRecordingService) GetRecordings(deviceID string, channelIdx int, begin, end time.Time, page int, pageSize int) ([]models.CloudRecording, error) {
	if deviceID == "" {
		return nil, errors.New("invalid request: device ID must be provided")
	}

	if begin.IsZero() || end.IsZero() || begin.After(end) {
		return nil, errors.New("invalid request: begin and end times must be valid and begin must be before end")
	}

	if channelIdx < 0 {
		return nil, errors.New("invalid request: channel index must be a non-negative integer")
	}

	if page <= 0 {
		return nil, errors.New("invalid request: page must be a positive integer")
	}
	if pageSize <= 0 || pageSize > 1000 {
		return nil, errors.New("invalid request: page size must be a positive integer and less than or equal to 1000")
	}
	limit := 1000
	offset := (page - 1) * pageSize

	recordings, err := s.cloudRecordingModels.List(deviceID, channelIdx, begin, end, limit, offset)
	if err != nil {
		return nil, err
	}

	return recordings, nil
}
