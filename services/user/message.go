package user

import (
	"encoding/json"
	"nexus/models"
	"time"
)

type Resource struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type AlarmContent struct {
	Type    string  `json:"type"`
	Resource Resource `json:"resource"`
	Device  string  `json:"device"`
}

type DeviceNotificationContent struct {
	Action string `json:"action"`
	Device string `json:"device"`
}

type IMessage interface {
	GetType() string
}

type BaseMessage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (message *BaseMessage) GetType() string {
	return message.Type
}

type AlarmMessage struct {
	BaseMessage
	Content AlarmContent `json:"content"`
}

type DeviceNotificationMessage struct {
	BaseMessage
	Content DeviceNotificationContent `json:"content"`
}

type IMessageService interface {
	ListByUserAndType(userID *string, messageType *string, page int, pageSize int) ([]string, error)
	Create(userID string, messageType string, content string) (string, error)
	Update(id string) error
}

type MessageService struct {
	messageModels models.IMessageModels
}

func NewMessageService() *MessageService {
	return &MessageService{
		messageModels: &models.MessageModels{},
	}
}

func NewMessageServiceWithModel(messageModels models.IMessageModels) *MessageService {
	return &MessageService{
		messageModels,
	}
}

func (service *MessageService) ListByUserAndType(userID *string, messageType *string, page int, pageSize int) ([]IMessage, error) {
	messages, err := service.messageModels.List(userID, messageType, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	var result = toMessages(messages)
	return result, nil
}

func (service *MessageService) Create(userID string, messageType string, content string) (string, error) {
	message := models.Message{
		ID:      models.GenerateID(),
		UserID:  userID,
		Type:    messageType,
		Content: []byte(content),
	}
	err := service.messageModels.Create(&message)
	if err != nil {
		return "", err
	}
	return message.ID, nil
}

func (service *MessageService) Delete(id string) error {
	message, err := service.messageModels.Get(id)
	if err != nil {
		return err
	}
	message.Status = false
	message.UpdatedAt = time.Now()
	err = service.messageModels.Update(&message)
	if err != nil {
		return err
	}
	return nil
}

/* internal function */

func toMessages(messages []models.Message) []IMessage {
	var result []IMessage
	for _, message := range messages {
		switch message.Type {
		case "alarm":
			alarmMessage, err := toAlarmMessage(message)
			if err != nil {
				continue
			}
			result = append(result, alarmMessage)
		case "device_notification":
			deviceNotificationMessage, err := toDeviceNotificationMessage(message)
			if err != nil {
				continue
			}
			result = append(result, deviceNotificationMessage)
		default:
			// Handle unknown message type
			// You can log an error or return a default message
			// For now, we will just skip it
			continue
		}
	}
	return result
}

func toAlarmMessage(message models.Message) (*AlarmMessage, error) {
	var AlarmContent AlarmContent
	if err := json.Unmarshal(message.Content, &AlarmContent); err != nil {
		return nil, err
	}
	return &AlarmMessage{
		BaseMessage: BaseMessage{
			ID:        message.ID,
			UserID:    message.UserID,
			Type:      message.Type,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		},
		Content: AlarmContent,
	}, nil
}

func toDeviceNotificationMessage(message models.Message) (*DeviceNotificationMessage, error) {
	var DeviceNotificationContent DeviceNotificationContent
	if err := json.Unmarshal(message.Content, &DeviceNotificationContent); err != nil {
		return nil, err
	}
	return &DeviceNotificationMessage{
		BaseMessage: BaseMessage{
			ID:        message.ID,
			UserID:    message.UserID,
			Type:      message.Type,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		},
		Content: DeviceNotificationContent,
	}, nil
}
