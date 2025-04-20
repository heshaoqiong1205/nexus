package models

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID               string `gorm:"primaryKey"`
	Name             string
	Description      string
	RequiredFeatures datatypes.JSON
	CreatedAt        time.Time
	UpdatedAt        time.Time `gorm:"autoUpdateTime:false"`
	Status           bool
}

type IProductModels interface {
	Get(id string) (Product, error)
	Create(product *Product) error
	Update(product *Product) error
}

type ProductModels struct {
}

func (models *ProductModels) Get(id string) (Product, error) {
	var product Product
	result := db.Find(&product, "id = ?", id)
	if result.Error != nil {
		return Product{}, result.Error
	}
	return product, nil
}

func (models *ProductModels) Create(product *Product) error {
	result := db.Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *ProductModels) Update(product *Product) error {
	result := db.Save(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
