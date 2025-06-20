package services_testing

import (
	"errors"
	external_storage "nexus/external/storage"
	"nexus/models"
	"nexus/services/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCloudStorageModels struct {
	mock.Mock
}

func (m *MockCloudStorageModels) Get(id string) (models.CloudStorage, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return models.CloudStorage{}, args.Error(1)
	}
	return args.Get(0).(models.CloudStorage), args.Error(1)
}

func (m *MockCloudStorageModels) GetByDeviceID(deviceIDs string) ([]models.CloudStorage, error) {
	args := m.Called(deviceIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CloudStorage), args.Error(1)
}

func (m *MockCloudStorageModels) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCloudStorageModels) List(limit int, offset int, orderBY string) ([]models.CloudStorage, error) {
	args := m.Called(limit, offset, orderBY)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CloudStorage), args.Error(1)
}

func (m *MockCloudStorageModels) Create(storage *models.CloudStorage) error {
	args := m.Called(storage)
	return args.Error(0)
}

func (m *MockCloudStorageModels) Update(storage *models.CloudStorage) error {
	args := m.Called(storage)
	return args.Error(0)
}

func (m *MockCloudStorageModels) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockBucketService struct {
	mock.Mock
}

func (m *MockBucketService) GetBucket(id string) (storage.BucketConfig, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return storage.BucketConfig{}, args.Error(1)
	}
	return args.Get(0).(storage.BucketConfig), args.Error(1)
}

func (m *MockBucketService) GetBucketByName(name string) (storage.BucketConfig, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return storage.BucketConfig{}, args.Error(1)
	}
	return args.Get(0).(storage.BucketConfig), args.Error(1)
}

func (m *MockBucketService) GetBucketByID(id string) (storage.BucketConfig, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return storage.BucketConfig{}, args.Error(1)
	}
	return args.Get(0).(storage.BucketConfig), args.Error(1)
}

func (m *MockBucketService) GetPaginatedBuckets(page int, pageSize int, orderBy string) (storage.PaginatedBuckets, error) {
	args := m.Called(page, pageSize, orderBy)
	if args.Get(0) == nil {
		return storage.PaginatedBuckets{}, args.Error(1)
	}
	return args.Get(0).(storage.PaginatedBuckets), args.Error(1)
}

func (m *MockBucketService) CreateBucket(bucket storage.BucketConfig) error {
	args := m.Called(bucket)
	return args.Error(0)
}

func (m *MockBucketService) UpdateBucket(bucket storage.BucketConfig) error {
	args := m.Called(bucket)
	return args.Error(0)
}
func (m *MockBucketService) DeleteBucket(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockStsService struct {
	mock.Mock
}

func (m *MockStsService) GenerateCredentials(provider string, region string, statement []external_storage.Statement, duaration int32) (external_storage.Credentials, error) {
	args := m.Called(provider, region, statement, duaration)
	if args.Get(0) == nil {
		return external_storage.Credentials{}, args.Error(1)
	}
	return args.Get(0).(external_storage.Credentials), args.Error(1)
}

// TestGetCloudStorage tests the GetCloudStorage method of the CloudStorageService

func TestGetCloudStorage(t *testing.T) {
	// Test the GetCloudStorage method
	bucketService := new(MockBucketService)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := external_storage.NewStsService()
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)
	mockCloudStorageModels.On("Get", "test-storage").Return(models.CloudStorage{
		ID:       "test-storage",
		DeviceID: "device-1",
		Tos:      "tos-1",
		Bucket:   "test-bucket",
		Mode:     "active",
		Status:   true,
	}, nil).Once()
	storageConfig, err := service.GetCloudStorage("test-storage")

	if err != nil {
		t.Errorf("GetCloudStorage returned an error: %v", err)
	}

	assert.Equal(t, storageConfig.ID, "test-storage")
	assert.Equal(t, storageConfig.DeviceID, "device-1")
	assert.Equal(t, storageConfig.Tos, "tos-1")
	assert.Equal(t, storageConfig.Bucket, "test-bucket")
	assert.Equal(t, storageConfig.Mode, "active")
	assert.Equal(t, storageConfig.Status, true)
	mockCloudStorageModels.AssertExpectations(t)
}

func TestGetCloudStorageNotFound(t *testing.T) {
	// Test the GetCloudStorage method with an invalid ID
	mockBucketModels := new(MockBucketModels)
	bucketService := storage.NewBucketService(mockBucketModels)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := external_storage.NewStsService()
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)
	mockCloudStorageModels.On("Get", "invalid-storage").Return(models.CloudStorage{}, errors.New("not found")).Once()
	_, err := service.GetCloudStorage("invalid-storage")

	if err == nil {
		t.Error("GetCloudStorage should have returned an error for an invalid ID")
	}

	assert.Equal(t, err, errors.New("not found"))
	mockCloudStorageModels.AssertExpectations(t)
}

func TestGetCloudStorageByDeviceID(t *testing.T) {
	// Test the GetCloudStorageByDeviceID method
	mockBucketModels := new(MockBucketModels)
	bucketService := storage.NewBucketService(mockBucketModels)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := external_storage.NewStsService()
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)
	mockCloudStorageModels.On("GetByDeviceID", "device-1").Return([]models.CloudStorage{
		{
			ID:       "storage-1",
			DeviceID: "device-1",
			Tos:      "tos-1",
			Bucket:   "bucket-1",
			Mode:     "active",
			Status:   true,
		},
	}, nil).Once()
	storages, err := service.GetCloudStorageByDeviceID("device-1")

	if err != nil {
		t.Errorf("GetCloudStorageByDeviceID returned an error: %v", err)
	}

	assert.Len(t, storages, 1)
	assert.Equal(t, storages[0].ID, "storage-1")
	assert.Equal(t, storages[0].DeviceID, "device-1")
	assert.Equal(t, storages[0].Tos, "tos-1")
	assert.Equal(t, storages[0].Bucket, "bucket-1")
	assert.Equal(t, storages[0].Mode, "active")
	assert.Equal(t, storages[0].Status, true)
	mockCloudStorageModels.AssertExpectations(t)
}

func TestUpdateCloudStorage(t *testing.T) {
	// Test the UpdateCloudStorage method
	mockBucketModels := new(MockBucketModels)
	bucketService := storage.NewBucketService(mockBucketModels)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := external_storage.NewStsService()
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)
	storageConfig := &storage.CloudStorage{
		ID:       "existing-storage",
		DeviceID: "device-1",
		Tos:      "media7",
		Bucket:   "test-bucket",
		Mode:     "plaintext",
		Status:   true,
	}

	mockCloudStorageModels.On("Update", &models.CloudStorage{
		ID:       storageConfig.ID,
		DeviceID: storageConfig.DeviceID,
		Tos:      storageConfig.Tos,
		Bucket:   storageConfig.Bucket,
		Mode:     storageConfig.Mode,
		Status:   storageConfig.Status,
	}).Return(nil).Once()

	mockBucketModels.On("GetByName", storageConfig.Bucket).Return(models.Bucket{
		ID:        "test-bucket",
		Name:      storageConfig.Bucket,
		Region:    "us-west-1",
		CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Provider:  "aws",
		Endpoint:  "s3.us-west-1.amazonaws.com",
		Status:    true,
	}, nil).Once()

	err := service.UpdateCloudStorage(storageConfig)

	assert.Equal(t, err, nil)
	mockCloudStorageModels.AssertExpectations(t)
}

func TestGetStorageConfigs(t *testing.T) {
	// Test the GetStorageConfigs method

	bucketService := new(MockBucketService)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := external_storage.NewStsService()
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)

	mockCloudStorageModels.On("GetByDeviceID", "device-1").Return([]models.CloudStorage{
		{
			ID:       "storage-1",
			DeviceID: "device-1",
			Tos:      "media7",
			Bucket:   "bucket-1",
			Mode:     "plaintext",
			Status:   true,
		},
	}, nil).Once()

	bucketService.On("GetBucketByName", "bucket-1").Return(storage.BucketConfig{
		ID:       "bucket-1",
		Name:     "Bucket 1",
		Region:   "us-west-1",
		Provider: "aws",
		Endpoint: "s3.us-west-1.amazonaws.com",
	}, nil).Once()

	configs, err := service.GetStorageConfigs("device-1")

	if err != nil {
		t.Errorf("GetStorageConfigs returned an error: %v", err)
	}

	assert.Equal(t, configs[0].Tos, "media7")
	assert.Equal(t, configs[0].Bucket, "Bucket 1")
	assert.Equal(t, configs[0].Mode, "plaintext")
	assert.Equal(t, configs[0].Region, "us-west-1")
	assert.Equal(t, configs[0].Provider, "aws")
	assert.Equal(t, configs[0].Endpoint, "s3.us-west-1.amazonaws.com")
	mockCloudStorageModels.AssertExpectations(t)
}

func TestGetStorageCredentials(t *testing.T) {
	// Test the GetStorageCredentials method
	bucketService := new(MockBucketService)
	mockCloudStorageModels := new(MockCloudStorageModels)
	stsService := new(MockStsService)
	service := storage.NewCloudStorageService(mockCloudStorageModels, bucketService, stsService)

	mockCloudStorageModels.On("GetByDeviceID", "device-1").Return([]models.CloudStorage{
		{
			ID:       "storage-1",
			DeviceID: "device-1",
			Tos:      "media7",
			Bucket:   "bucket-1",
			Mode:     "plaintext",
			Status:   true,
		},
	}, nil).Once()

	bucketService.On("GetBucketByName", "bucket-1").Return(storage.BucketConfig{
		ID:       "bucket-1",
		Name:     "Bucket 1",
		Region:   "us-west-1",
		Provider: "aws",
		Endpoint: "s3.us-west-1.amazonaws.com",
	}, nil).Once()

	accessKeyId := "test-access-key"
	secretAccessKey := "test-secret-key"
	sessionToken := "test-session-token"
	expiratioen := time.Now().Add(time.Hour)
	stsService.On("GenerateCredentials", "aws", "us-west-1", mock.Anything, int32(3600)).Return(external_storage.Credentials{
		AccessKeyId:     &accessKeyId,
		SecretAccessKey: &secretAccessKey,
		Expiration:      &expiratioen,
		SessionToken:    &sessionToken,
	}, nil).Once()
	storageConfig, err := service.GetStorageCredentials([]external_storage.Action{external_storage.GetObject, external_storage.PutObject}, "device-1")

	if err != nil {
		t.Errorf("GetStorageCredentials returned an error: %v", err)
	}

	assert.Equal(t, storageConfig[0].Bucket, "Bucket 1")
	assert.Equal(t, *storageConfig[0].Credentials.AccessKeyId, "test-access-key")
	assert.Equal(t, *storageConfig[0].Credentials.SecretAccessKey, "test-secret-key")
	assert.Equal(t, *storageConfig[0].Credentials.Expiration, expiratioen)
	assert.Equal(t, *storageConfig[0].Credentials.SessionToken, "test-session-token")
	mockCloudStorageModels.AssertExpectations(t)
}
