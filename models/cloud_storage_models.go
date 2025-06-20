package models

import "time"

type CloudStorage struct {
	ID        string `gorm:"primaryKey"`
	DeviceID  string
	Tos       string // Type of storage, e.g., "media7", "media14", "media30", "media180", "log", etc.
	Bucket    string
	Mode      string
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
	Status    bool
}

type ICloudStorageModels interface {
	Get(id string) (CloudStorage, error)
	GetByDeviceID(deviceID string) ([]CloudStorage, error)
	Count() (int64, error)
	List(limt int, offset int, orderBY string) ([]CloudStorage, error)
	Create(storage *CloudStorage) error
	Update(storage *CloudStorage) error
	Delete(id string) error
}

type CloudStorageModels struct {
}

func (models *CloudStorageModels) Get(id string) (CloudStorage, error) {
	var storage CloudStorage
	result := db.Find(&storage, "id = ?", id)
	if result.Error != nil {
		return CloudStorage{}, result.Error
	}
	return storage, nil
}

func (models *CloudStorageModels) GetByDeviceID(deviceIDs string) ([]CloudStorage, error) {
	var storages []CloudStorage

	result := db.Find(&storages, "device_id = ? and status = true", deviceIDs)
	if result.Error != nil {
		return nil, result.Error
	}
	return storages, nil
}

func (models *CloudStorageModels) Count() (int64, error) {
	var total int64

	result := db.Model(&CloudStorage{}).Where("status = true").Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *CloudStorageModels) List(limt int, offset int, orderBY string) ([]CloudStorage, error) {
	var storages []CloudStorage

	result := db.Limit(limt).Offset(offset).Order(orderBY).Find(&storages, "status = true")
	if result.Error != nil {
		return nil, result.Error
	}
	return storages, nil
}

func (models *CloudStorageModels) Create(storage *CloudStorage) error {
	result := db.Create(storage)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudStorageModels) Update(storage *CloudStorage) error {
	result := db.Save(storage)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *CloudStorageModels) Delete(id string) error {
	result := db.Delete(&CloudStorage{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
