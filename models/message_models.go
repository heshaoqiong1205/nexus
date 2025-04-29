package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Message struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	Type      string
	Content   datatypes.JSON
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
	Status    bool
}

type IMessageModels interface {
	Create(message *Message) error
	Update(message *Message) error
	Get(id string) (Message, error)
	List(userID *string, messageType *string, limt int, offset int) ([]Message, error)
}

type MessageModels struct {
}

func (models *MessageModels) Create(message *Message) error {
	result := db.Create(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *MessageModels) Update(message *Message) error {
	result := db.Save(message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (models *MessageModels) Get(id string) (Message, error) {
	var message Message
	result := db.Find(&message, "id = ?", id)
	if result.Error != nil {
		return Message{}, result.Error
	}
	return message, nil
}

func (models *MessageModels) List(userID *string, messageType *string, limt int, offset int) ([]Message, error) {
	var messages []Message
	var result *gorm.DB

	if userID != nil && messageType != nil {
		result = db.Limit(limt).Offset(offset).Order("created desc").Find(&messages, "user_id = ? and type = ?", userID, messageType)
	} else if userID != nil {
		result = db.Limit(limt).Offset(offset).Order("created desc").Find(&messages, "user_id = ?", userID)
	} else if messageType != nil {
		result = db.Limit(limt).Offset(offset).Order("created desc").Find(&messages, "type = ?", messageType)
	} else {
		result = db.Limit(limt).Offset(offset).Order("created desc").Find(&messages)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}
