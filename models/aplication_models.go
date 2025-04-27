package models

import "time"

type Application struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	SecretKey   string
	Salt       string
	CreatedAt   time.Time
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
	Status      bool
}

type IApplicationModels interface {
	Get(id string) (Application, error)
	List(limt int, offset int, orderBY string) ([]Application, error)
	Create(application *Application) error
	Update(application *Application) error
	Count() (int64, error)
}

type ApplicationModels struct {
}

func (models *ApplicationModels) Get(id string) (Application, error) {
	var application Application
	result := db.Find(&application, "id = ?", id)
	if result.Error != nil {
		return Application{}, result.Error
	}
	return application, nil
}

func (models *ApplicationModels) Count() (int64, error) {
	var total int64

	result := db.Model(&Application{}).Where("status = true").Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *ApplicationModels) List(limt int, offset int, orderBY string) ([]Application, error) {
	var applications []Application

	result := db.Limit(limt).Offset(offset).Order(orderBY).Find(&applications, "status = true")
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (models *ApplicationModels) Create(application *Application) error {
	result := db.Create(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *ApplicationModels) Update(application *Application) error {
	result := db.Save(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
