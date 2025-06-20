package storage

import "nexus/models"

type BucketConfig struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Provider string `json:"provider"`
}

type PaginatedBuckets struct {
	Total   int64          `json:"total"`
	Buckets []BucketConfig `json:"buckets"`
}

type IBucketService interface {
	GetBucket(ID string) (BucketConfig, error)
	GetBucketByName(bucketName string) (BucketConfig, error)
	GetPaginatedBuckets(page int, pageSize int, orderBy string) (PaginatedBuckets, error)
	CreateBucket(bucketConfig BucketConfig) error
	UpdateBucket(bucketConfig BucketConfig) error
	DeleteBucket(bucketID string) error
}

type BucketService struct {
	bucketModels models.IBucketModels
}

func NewBucketService(bucketModels models.IBucketModels) IBucketService {
	return &BucketService{
		bucketModels: bucketModels,
	}
}

func (s *BucketService) GetBucket(ID string) (BucketConfig, error) {
	bucket, err := s.bucketModels.Get(ID)
	if err != nil {
		return BucketConfig{}, err
	}

	return BucketConfig{
		ID:       bucket.ID,
		Name:     bucket.Name,
		Region:   bucket.Region,
		Endpoint: bucket.Endpoint,
		Provider: bucket.Provider,
	}, nil
}

func (s *BucketService) GetBucketByName(name string) (BucketConfig, error) {
	bucket, err := s.bucketModels.GetByName(name)
	if err != nil {
		return BucketConfig{}, err
	}

	return BucketConfig{
		ID:       bucket.ID,
		Name:     bucket.Name,
		Region:   bucket.Region,
		Endpoint: bucket.Endpoint,
		Provider: bucket.Provider,
	}, nil
}

func (s *BucketService) GetPaginatedBuckets(page int, pageSize int, orderBy string) (PaginatedBuckets, error) {
	if page < 1 {
		page = 1
	}
	limit := pageSize
	offset := (page - 1) * pageSize

	buckets, err := s.bucketModels.List(limit, offset, orderBy)
	if err != nil {
		return PaginatedBuckets{}, err
	}

	total, err := s.bucketModels.Count()
	if err != nil {
		return PaginatedBuckets{}, err
	}

	var bucketConfigs []BucketConfig
	for _, bucket := range buckets {
		bucketConfigs = append(bucketConfigs, BucketConfig{
			ID:       bucket.ID,
			Name:     bucket.Name,
			Region:   bucket.Region,
			Endpoint: bucket.Endpoint,
			Provider: bucket.Provider,
		})
	}

	return PaginatedBuckets{Total: total, Buckets: bucketConfigs}, nil
}

func (s *BucketService) CreateBucket(bucketConfig BucketConfig) error {
	bucket := models.Bucket{
		ID:       bucketConfig.ID,
		Name:     bucketConfig.Name,
		Region:   bucketConfig.Region,
		Endpoint: bucketConfig.Endpoint,
		Provider: bucketConfig.Provider,
		Status:   true, // Assuming new buckets are active by default
	}

	return s.bucketModels.Create(&bucket)
}

func (s *BucketService) UpdateBucket(bucketConfig BucketConfig) error {
	bucket, err := s.bucketModels.Get(bucketConfig.ID)
	if err != nil {
		return err
	}

	bucket.Name = bucketConfig.Name
	bucket.Region = bucketConfig.Region
	bucket.Endpoint = bucketConfig.Endpoint
	bucket.Provider = bucketConfig.Provider

	return s.bucketModels.Update(&bucket)
}

func (s *BucketService) DeleteBucket(bucketID string) error {
	return s.bucketModels.Delete(bucketID)
}
