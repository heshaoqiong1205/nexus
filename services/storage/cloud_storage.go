package storage

import (
	"errors"
	external_storage "nexus/external/storage"
	"nexus/models"
)

const (
	Media7   = "media7"
	Media14  = "media14"
	Media30  = "media30"
	Media180 = "media180"
	Log      = "log"
)

const (
	Plaintext = "plaintext"
	Encrypted = "encrypted"
)

type StorageConfig struct {
	Tos      string `json:"tos"`
	Bucket   string `json:"bucket"`
	Mode     string `json:"mode"`
	Path     string `json:"path"`
	Region   string `json:"region"`
	Provider string `json:"provider"`
	Endpoint string `json:"endpoint"`
}

type CloudStorage struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
	Tos      string `json:"tos"` // Type of storage, e.g., "media7", "media14", "media30", "media180", "log", etc.
	Bucket   string `json:"bucket"`
	Mode     string `json:"mode"`
	Path     string `json:"path"`
	Status   bool   `json:"status"`
}

type StorageCredentials struct {
	Bucket      string                       `json:"bucket"`
	Credentials external_storage.Credentials `json:"credentials"`
}

type ICloudStorageService interface {
	GetCloudStorage(id string) (CloudStorage, error)
	GetCloudStorageByDeviceID(deviceID string) ([]CloudStorage, error)
	UpdateCloudStorage(storage *CloudStorage) error
	GetStorageConfigs(deviceID string) ([]StorageConfig, error)
	GetStorageCredentials(actions []external_storage.Action, deviceID string) ([]StorageCredentials, error)
}

type CloudStorageService struct {
	cloudStorageModels models.ICloudStorageModels
	bucketService      IBucketService
	stsService         external_storage.IStsService
}

func NewCloudStorageService() *CloudStorageService {
	return &CloudStorageService{
		cloudStorageModels: &models.CloudStorageModels{},
		bucketService:      NewBucketService(),
		stsService:         external_storage.NewStsService(),
	}
}


func NewCloudStorageServiceWithModels(cloudStorageModels models.ICloudStorageModels, bucketService IBucketService, stsService external_storage.IStsService) *CloudStorageService {
	return &CloudStorageService{
		cloudStorageModels,
		bucketService,
		stsService,
	}
}

func (cs *CloudStorageService) GetCloudStorage(id string) (CloudStorage, error) {
	cloudStorage, err := cs.cloudStorageModels.Get(id)
	if err != nil {
		return CloudStorage{}, err
	}

	return CloudStorage{
		ID:       cloudStorage.ID,
		DeviceID: cloudStorage.DeviceID,
		Tos:      cloudStorage.Tos,
		Bucket:   cloudStorage.Bucket,
		Mode:     cloudStorage.Mode,
		Path:     cloudStorage.Path,
		Status:   cloudStorage.Status,
	}, nil
}

func (cs *CloudStorageService) GetCloudStorageByDeviceID(deviceID string) ([]CloudStorage, error) {
	cloudStorages, err := cs.cloudStorageModels.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}
	var result []CloudStorage
	for _, storage := range cloudStorages {
		result = append(result, CloudStorage{
			ID:       storage.ID,
			DeviceID: storage.DeviceID,
			Tos:      storage.Tos,
			Bucket:   storage.Bucket,
			Mode:     storage.Mode,
			Path:     storage.Path,
			Status:   storage.Status,
		})
	}
	return result, nil
}

func (cs *CloudStorageService) UpdateCloudStorage(storage *CloudStorage) error {
	if storage == nil {
		return errors.New("storage cannot be nil")
	}
	if storage.ID == "" {
		return errors.New("storage ID cannot be empty")
	}
	if storage.Mode != Plaintext && storage.Mode != Encrypted {
		return errors.New("storage mode must be either 'plaintext' or 'encrypted'")
	}
	if storage.Tos != Media7 && storage.Tos != Media14 && storage.Tos != Media30 && storage.Tos != Media180 && storage.Tos != Log {
		return errors.New("invalid type of storage")
	}
	_, err := cs.bucketService.GetBucketByName(storage.Bucket)
	if err != nil {
		return err
	}

	return cs.cloudStorageModels.Update(&models.CloudStorage{
		ID:       storage.ID,
		DeviceID: storage.DeviceID,
		Tos:      storage.Tos,
		Bucket:   storage.Bucket,
		Mode:     storage.Mode,
		Status:   storage.Status,
	})
}

func (cs *CloudStorageService) GetStorageConfigs(deviceID string) ([]StorageConfig, error) {
	cloudStorages, err := cs.cloudStorageModels.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	var storageConfigs []StorageConfig
	for _, storage := range cloudStorages {
		bucket, err := cs.bucketService.GetBucketByName(storage.Bucket)
		if err != nil {
			return nil, err
		}
		storageConfigs = append(storageConfigs, StorageConfig{
			Tos:      storage.Tos,
			Bucket:   bucket.Name,
			Mode:     storage.Mode,
			Path:     storage.Path,
			Region:   bucket.Region,
			Provider: bucket.Provider,
			Endpoint: bucket.Endpoint,
		})

	}
	return storageConfigs, nil
}

func (cs *CloudStorageService) GetStorageCredentials(actions []external_storage.Action, deviceID string) ([]StorageCredentials, error) {
	if len(actions) == 0 {
		return nil, errors.New("actions cannot be empty")
	}
	cloudStorages, err := cs.cloudStorageModels.GetByDeviceID(deviceID)
	if err != nil {
		return nil, err
	}

	bucketMap := make(map[string][]string)
	for _, cloudStorage := range cloudStorages {
		bucket, exists := bucketMap[cloudStorage.Bucket]
		if !exists {
			bucketMap[cloudStorage.Bucket] = []string{cloudStorage.Path}
		}
		bucketMap[cloudStorage.Bucket] = append(bucket, cloudStorage.Path)
	}

	var storageCredentials []StorageCredentials
	for bucketName, paths := range bucketMap {
		bucket, err := cs.bucketService.GetBucketByName(bucketName)
		if err != nil {
			return nil, err
		}
		statements := make([]external_storage.Statement, len(paths))
		for _, path := range paths {
			statement := &external_storage.Statement{
				Effect:  external_storage.Allow,
				Actions: actions,
				Resource: external_storage.Resource{
					Bucket: bucket.Name,
					Path:   path,
				},
			}
			statements = append(statements, *statement)
		}

		credentials, err := cs.stsService.GenerateCredentials(bucket.Provider, bucket.Region, statements, 3600)
		if err != nil {
			return nil, err
		}
		storageCredentials = append(storageCredentials, StorageCredentials{
			Bucket:      bucket.Name,
			Credentials: credentials,
		})
	}
	return storageCredentials, nil
}
