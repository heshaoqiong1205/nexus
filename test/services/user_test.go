package services_testing

import (
	"nexus/models"
	"nexus/services/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	 ================================================================
				user service test

===================================================================
*/
type MockUserModels struct {
	// Mock methods for user models
}

func mockUser0() models.User {
	return models.User{
		ID:       "test_id0",
		Account:  "test_account0",
		Username: "test_username0",
		Password: "test_password0",
		Region:   "china",
		Location: "test_location0",
		Icon:     "test_icon0",
		Role:     "admin",
		Status:   true,
	}
}

func mockUser1() models.User {
	return models.User{
		ID:       "test_id1",
		Account:  "test_account1",
		Username: "test_username1",
		Password: "test_password1",
		Region:   "china",
		Location: "test_location1",
		Icon:     "test_icon1",
		Role:     "user",
		Status:   true,
	}
}

func (m *MockUserModels) Get(id string) (models.User, error) {
	// Mock implementation
	return mockUser0(), nil
}

func (m *MockUserModels) GetByAccount(account string) (models.User, error) {
	// Mock implementation
	return mockUser0(), nil
}

func (m *MockUserModels) Create(user *models.User) error {
	// Mock implementation
	return nil
}

func (m *MockUserModels) Update(user *models.User) error {
	// Mock implementation
	return nil
}

func (m *MockUserModels) List(LikeAccount string, limt int, offset int, orderBY string) ([]models.User, error) {
	// Mock implementation
	return []models.User{mockUser0(), mockUser1()}, nil
}

func TestSignUp(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the SignUp method
	userDetails, err := userService.SignUp("test_account0", "test_password0", "test_username0", "china")
	assert.Equal(t, nil, err)

	// Check the returned user details
	expected := &user.UserDetails{
		Account:       "test_account0",
		Username:      "test_username0",
		Region:        "china",
		Location:      "",
		Icon:          "",
		Role:          "user",
		LastLoginTime: userDetails.LastLoginTime,
		SginUpAt:      userDetails.SginUpAt,
	}
	assert.Equal(t, expected, userDetails)
}

func TestSignIn(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the SignIn method
	_, err := userService.SignIn("test_account0", "test_password0")
	assert.Equal(t, nil, err)

}

func TestGetUser(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the Get method
	userDetails, err := userService.Get("test_id0")
	assert.Equal(t, nil, err)

	// Check the returned user details
	expected := toUserDetails(mockUser0())
	assert.Equal(t, expected, userDetails)
}

func TestGetByAccount(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the GetByAccount method
	userDetails, err := userService.GetByAccount("test_account0")
	assert.Equal(t, nil, err)

	// Check the returned user details
	expected := toUserDetails(mockUser0())
	assert.Equal(t, expected, userDetails)
}

func TestSearchByAccount(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the List method
	users, err := userService.SearchByAccount("test_account", 10, 0, "account")
	assert.Equal(t, nil, err)

	// Check the returned users
	expected := []user.UserDetails{
		*toUserDetails(mockUser0()),
		*toUserDetails(mockUser1()),
	}
	assert.Equal(t, expected, users)
}

func TestUpdateUser(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	location := "test_location0"

	request := user.ModifyUserRequest{
		ID:       "test_id0",
		Username: nil,
		Location: &location,
		Icon:     nil,
	}
	// Call the Update method
	err := userService.Update(request)
	assert.Equal(t, nil, err)
}

/* internal function */
func toUserDetails(u models.User) *user.UserDetails {
	return &user.UserDetails{
		Account:       u.Account,
		Username:      u.Username,
		Region:        u.Region,
		Location:      u.Location,
		Icon:          u.Icon,
		Role:          u.Role,
		LastLoginTime: u.LastLoginTime,
		SginUpAt:      u.CreatedAt,
	}
}

/* ==================================================================
			message service test
=================================================================== */

type MockMessageModels struct {
	// Mock methods for message models
}

func mockMessage0() models.Message {
	return models.Message{
		ID:        "test_id0",
		UserID:    "test_user_id",
		Type:      "alarm",
		Content:   []byte("{\"type\":\"alarm\",\"resouce\":{\"type\":\"image\",\"url\":\"test_url\"},\"device\":\"test_device0\"}"),
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func mockMessage1() models.Message {
	return models.Message{
		ID:        "test_id1",
		UserID:    "test_user_id",
		Type:      "device_notification",
		Content:   []byte("{\"action\":\"ADD_DEVICE\",\"device\":\"test_device1\"}"),
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func (m *MockMessageModels) Get(id string) (models.Message, error) {
	// Mock implementation
	return mockMessage0(), nil
}

func (m *MockMessageModels) Create(message *models.Message) error {
	// Mock implementation
	return nil
}

func (m *MockMessageModels) Update(message *models.Message) error {
	// Mock implementation
	return nil
}

func (m *MockMessageModels) List(userID *string, messageType *string, limt int, offset int) ([]models.Message, error) {
	// Mock implementation
	return []models.Message{mockMessage0(), mockMessage1()}, nil
}

func TestCreateMessage(t *testing.T) {
	// Mock the message models
	mockMessageModels := &MockMessageModels{}

	// Create a new message service with the mock message models
	messageService := user.NewMessageServiceWithModel(mockMessageModels)

	// Call the Create method
	messageID, err := messageService.Create("test_user_id", "alarm", "test_content0")
	assert.Equal(t, nil, err)

	// Check the returned message ID
	assert.NotEqual(t, "", messageID)
}

func TestUpdateMessage(t *testing.T) {
	// Mock the message models
	mockMessageModels := &MockMessageModels{}

	// Create a new message service with the mock message models
	messageService := user.NewMessageServiceWithModel(mockMessageModels)

	// Call the Update method
	err := messageService.Delete("test_id0")
	assert.Equal(t, nil, err)
}

func TestMessageListByUserAndType(t *testing.T) {
	// Mock the message models
	mockMessageModels := &MockMessageModels{}

	// Create a new message service with the mock message models
	messageService := user.NewMessageServiceWithModel(mockMessageModels)

	// Call the ListByUserAndType method
	messages, err := messageService.ListByUserAndType(nil, nil, 1, 10)
	assert.Equal(t, nil, err)

	message0 := mockMessage0()
	message1 := mockMessage1()
	// Check the returned messages
	expected := []user.IMessage{
		&user.AlarmMessage{
			BaseMessage: user.BaseMessage{
				ID:        message0.ID,
				UserID:    message0.UserID,
				Type:      message0.Type,
				CreatedAt: message0.CreatedAt,
				UpdatedAt: message0.UpdatedAt,
			},
			Conntent: user.AlarmContent{
				Type: "alarm",
				Resouce: user.Resouce{
					Type: "image",
					Url:  "test_url",
				},
				Device: "test_device0",
			},
		},
		&user.DeviceNotificationMessage{
			BaseMessage: user.BaseMessage{
				ID:        message1.ID,
				UserID:    message1.UserID,
				Type:      message1.Type,
				CreatedAt: message1.CreatedAt,
				UpdatedAt: message1.UpdatedAt,
			},
			Conntent: user.DeviceNotificationContent{
				Action: "ADD_DEVICE",
				Device: "test_device1",
			},
		},
	}
	assert.Equal(t, 2, len(messages))
	assert.Equal(t, expected[0], messages[0])
	assert.Equal(t, expected[1], messages[1])
}
