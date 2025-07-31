package models

import "time"

type User struct {
	ID            string `gorm:"primaryKey"`
	Account       string
	Username      string
	Password      string
	Region        string
	Location      *Point `gorm:"type:point"`
	Icon          string
	Role          string
	LastLoginTime time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time `gorm:"autoUpdateTime:false"`
	Status        bool
}

type IUserModels interface {
	Get(id string) (User, error)
	GetByAccount(account string) (User, error)
	Count(likeAccount string) (int64, error)
	List(likeAccount string, limit int, offset int, orderBy string) ([]User, error)
	Create(user *User) error
	Update(user *User) error
}

type UserModels struct {
}

func (models *UserModels) Get(id string) (User, error) {
	var user User
	result := db.Find(&user, "id = ?", id)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (models *UserModels) GetByAccount(account string) (User, error) {
	var user User
	result := db.Find(&user, "account = ?", account)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (models *UserModels) Count(likeAccount string) (int64, error) {
	var total int64

	result := db.Model(&User{}).Where("status = true and account like ?", likeAccount).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (models *UserModels) List(likeAccount string, limit int, offset int, orderBy string) ([]User, error) {
	var users []User

	result := db.Limit(limit).Offset(offset).Order(orderBy).Find(&users, "status = true and account like ?", likeAccount)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (models *UserModels) Create(user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *UserModels) Update(user *User) error {
	result := db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

