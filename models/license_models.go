package models

import "time"

type License struct {
	ID        string `gorm:"primaryKey"`
	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    bool
}

type ILicenseModels interface {
	Get(id string) (License, error)
	Create(license *License) error
	Update(license *License) error
}

type LicenseModels struct {
}

func (models *LicenseModels) Get(id string) (License, error) {
	var license License
	result := db.Find(&license, "id = ?", id)
	if result.Error != nil {
		return License{}, result.Error
	}
	return license, nil
}

func (models *LicenseModels) Create(license *License) error {
	result := db.Create(license)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *LicenseModels) Update(license *License) error {
	result := db.Save(license)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
