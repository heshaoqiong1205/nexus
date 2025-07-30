package services_testing

import (
	"errors"
	"nexus/models"
	"nexus/services/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBucketModels struct {
	mock.Mock
}

func (m *MockBucketModels) Get(id string) (models.Bucket, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return models.Bucket{}, args.Error(1)
	}
	return args.Get(0).(models.Bucket), args.Error(1)
}

func (m *MockBucketModels) GetByName(name string) (models.Bucket, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return models.Bucket{}, args.Error(1)
	}
	return args.Get(0).(models.Bucket), args.Error(1)
}

func (m *MockBucketModels) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBucketModels) List(limit int, offset int, orderBY string) ([]models.Bucket, error) {
	args := m.Called(limit, offset, orderBY)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Bucket), args.Error(1)
}

func (m *MockBucketModels) Create(bucket *models.Bucket) error {
	args := m.Called(bucket)
	return args.Error(0)
}

func (m *MockBucketModels) Update(bucket *models.Bucket) error {
	args := m.Called(bucket)
	return args.Error(0)
}

func (m *MockBucketModels) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetBucket(t *testing.T) {
	// Test the GetBucket method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	mockBucketModels.On("Get", "test-bucket").Return(models.Bucket{
		ID:        "test-bucket",
		Name:      "Test Bucket",
		Region:    "us-west-1",
		CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Provider:  "aws",
		Endpoint:  "s3.us-west-1.amazonaws.com",
		Status:    true,
	}, nil).Once()
	config, err := service.GetBucket("test-bucket")

	if err != nil {
		t.Errorf("GetBucket returned an error: %v", err)
	}

	assert.Equal(t, config.ID, "test-bucket")
	assert.Equal(t, config.Name, "Test Bucket")
	assert.Equal(t, config.Region, "us-west-1")
	assert.Equal(t, config.Provider, "aws")
	assert.Equal(t, config.Endpoint, "s3.us-west-1.amazonaws.com")
	mockBucketModels.AssertExpectations(t)
	// Test the GetBucket method with an invalid bucket name

}

func TestGetBucketNotFound(t *testing.T) {
	// Test the GetBucket method with an invalid bucket name
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	mockBucketModels.On("Get", "invalid-bucket").Return(models.Bucket{}, errors.New("not found")).Once()
	_, err := service.GetBucket("invalid-bucket")

	if err == nil {
		t.Error("GetBucket should have returned an error for an invalid bucket name")
	}

	assert.Equal(t, err, errors.New("not found"))
	mockBucketModels.AssertExpectations(t)
}

func TestGetBucketByName(t *testing.T) {
	// Test the GetBucketByName method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	mockBucketModels.On("GetByName", "test-bucket").Return(models.Bucket{
		ID:        "test-bucket",
		Name:      "Test Bucket",
		Region:    "us-west-1",
		CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Provider:  "aws",
		Endpoint:  "s3.us-west-1.amazonaws.com",
		Status:    true,
	}, nil).Once()
	config, err := service.GetBucketByName("test-bucket")

	if err != nil {
		t.Errorf("GetBucketByName returned an error: %v", err)
	}

	assert.Equal(t, config.ID, "test-bucket")
	assert.Equal(t, config.Name, "Test Bucket")
	assert.Equal(t, config.Region, "us-west-1")
	assert.Equal(t, config.Provider, "aws")
	assert.Equal(t, config.Endpoint, "s3.us-west-1.amazonaws.com")
	mockBucketModels.AssertExpectations(t)
}

func TestListBuckets(t *testing.T) {
	// Test the ListBuckets method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	mockBucketModels.On("List", 10, 0, "name").Return([]models.Bucket{
		{
			ID:        "test-bucket-1",
			Name:      "Test Bucket 1",
			Region:    "us-west-1",
			CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			Provider:  "aws",
			Endpoint:  "s3.us-west-1.amazonaws.com",
			Status:    true,
		},
		{
			ID:        "test-bucket-2",
			Name:      "Test Bucket 2",
			Region:    "us-east-1",
			CreatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			Provider:  "aws",
			Endpoint:  "s3.us-east-1.amazonaws.com",
			Status:    true,
		},
	}, nil).Once()
	mockBucketModels.On("Count").Return(int64(2), nil).Once()
	paginated, err := service.GetPaginatedBuckets(1, 10, "name")

	if err != nil {
		t.Errorf("ListBuckets returned an error: %v", err)
	}

	assert.Len(t, paginated.Buckets, 2)
	assert.Equal(t, paginated.Total, int64(2))
	assert.Equal(t, paginated.Buckets[0].ID, "test-bucket-1")
	assert.Equal(t, paginated.Buckets[0].Name, "Test Bucket 1")
	assert.Equal(t, paginated.Buckets[0].Region, "us-west-1")
	assert.Equal(t, paginated.Buckets[0].Provider, "aws")
	assert.Equal(t, paginated.Buckets[0].Endpoint, "s3.us-west-1.amazonaws.com")
	assert.Equal(t, paginated.Buckets[1].ID, "test-bucket-2")
	assert.Equal(t, paginated.Buckets[1].Name, "Test Bucket 2")
	assert.Equal(t, paginated.Buckets[1].Region, "us-east-1")
	assert.Equal(t, paginated.Buckets[1].Provider, "aws")
	assert.Equal(t, paginated.Buckets[1].Endpoint, "s3.us-east-1.amazonaws.com")
	mockBucketModels.AssertExpectations(t)
}

func TestCreateBucket(t *testing.T) {
	// Test the CreateBucket method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	bucketConfig := storage.BucketConfig{
		ID:       "new-bucket",
		Name:     "New Bucket",
		Region:   "us-west-2",
		Endpoint: "s3.us-west-2.amazonaws.com",
		Provider: "aws",
	}

	mockBucketModels.On("Create", &models.Bucket{
		ID:       bucketConfig.ID,
		Name:     bucketConfig.Name,
		Region:   bucketConfig.Region,
		Endpoint: bucketConfig.Endpoint,
		Provider: bucketConfig.Provider,
		Status:   true, // Assuming new buckets are active by default
	}).Return(nil).Once()

	err := service.CreateBucket(bucketConfig)

	assert.Equal(t, err, nil)
	mockBucketModels.AssertExpectations(t)
}

func TestUpdateBucket(t *testing.T) {
	// Test the UpdateBucket method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	bucketConfig := storage.BucketConfig{
		ID:       "existing-bucket",
		Name:     "Updated Bucket",
		Region:   "us-west-2",
		Endpoint: "s3.us-west-2.amazonaws.com",
		Provider: "aws",
	}

	mockBucketModels.On("Get", bucketConfig.ID).Return(models.Bucket{
		ID:       bucketConfig.ID,
		Name:     "Old Bucket",
		Region:   "us-east-1",
		Endpoint: "s3.us-east-1.amazonaws.com",
		Provider: "aws",
		Status:   true, // Assuming buckets remain active after update
	}, nil).Once()

	mockBucketModels.On("Update", &models.Bucket{
		ID:       bucketConfig.ID,
		Name:     bucketConfig.Name,
		Region:   bucketConfig.Region,
		Endpoint: bucketConfig.Endpoint,
		Provider: bucketConfig.Provider,
		Status:   true, // Assuming buckets remain active after update
	}).Return(nil).Once()

	err := service.UpdateBucket(bucketConfig)

	assert.Equal(t, err, nil)
	mockBucketModels.AssertExpectations(t)
}

func TestDeleteBucket(t *testing.T) {
	// Test the DeleteBucket method
	mockBucketModels := new(MockBucketModels)
	service := storage.NewBucketServiceWithModels(mockBucketModels)
	bucketID := "bucket-to-delete"

	mockBucketModels.On("Delete", bucketID).Return(nil).Once()

	err := service.DeleteBucket(bucketID)

	assert.Equal(t, err, nil)
	mockBucketModels.AssertExpectations(t)
}
