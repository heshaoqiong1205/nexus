package models

import "time"

type Bucket struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Region    string
	Provider  string
	Endpoint  string
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
	Status    bool
}

type IBucketModels interface {
	Get(id string) (Bucket, error)
	GetByName(name string) (Bucket, error)
	Count() (int64, error)
	List(limt int, offset int, orderBY string) ([]Bucket, error)
	Create(bucket *Bucket) error
	Update(bucket *Bucket) error
	Delete(id string) error
}

type BucketModels struct {
}

func (models *BucketModels) Get(id string) (Bucket, error) {
	var bucket Bucket
	result := db.Find(&bucket, "id = ? and status = true", id)
	if result.Error != nil {
		return Bucket{}, result.Error
	}
	return bucket, nil
}

func (models *BucketModels) GetByName(name string) (Bucket, error) {
	var bucket Bucket
	result := db.Find(&bucket, "name = ? and status = true", name)
	if result.Error != nil {
		return Bucket{}, result.Error
	}
	return bucket, nil
}

func (models *BucketModels) Count() (int64, error) {
	var total int64

	result := db.Model(&Bucket{}).Where("status = true").Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *BucketModels) List(limt int, offset int, orderBY string) ([]Bucket, error) {
	var buckets []Bucket

	result := db.Limit(limt).Offset(offset).Order(orderBY).Find(&buckets, "status = true")
	if result.Error != nil {
		return nil, result.Error
	}
	return buckets, nil
}

func (models *BucketModels) Create(bucket *Bucket) error {
	result := db.Create(bucket)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *BucketModels) Update(bucket *Bucket) error {
	result := db.Save(bucket)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *BucketModels) Delete(id string) error {
	var bucket Bucket
	result := db.Delete(&bucket, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
