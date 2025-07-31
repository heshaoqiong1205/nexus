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
		Location:  &models.Point{Latitude: 39.9042, Longitude: 116.4074},
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
		Location: &models.Point{Latitude: 31.2304, Longitude: 121.4737},
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

func (m *MockUserModels) List(LikeAccount string, limit int, offset int, orderBY string) ([]models.User, error) {
	// Mock implementation
	return []models.User{mockUser0(), mockUser1()}, nil
}

func (m *MockUserModels) Count(LikeAccount string) (int64, error) {
	// Mock implementation
	return int64(2), nil
}

func TestSignUp(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	location := user.Point{Latitude: 39.9042, Longitude: 116.4074}

	// Call the SignUp method
	request := user.SignUpRequest{
		Account:  "test_account0",
		Password: "test_password0",
		Username: "test_username0",
		Region:   "china",
		Location: &location,
	}
	signUpResponse, err := userService.SignUp(request)
	assert.Equal(t, nil, err)

	// Check the returned user details
	expected := user.SignUpResponse{
		Token: signUpResponse.Token,
		User: &user.UserDetails{
			Account:       "test_account0",
			Username:      "test_username0",
			Region:        "china",
			Location:      &location,
			Icon:          "",
			Role:          "user",
			LastLoginTime: signUpResponse.User.LastLoginTime,
			SignUpAt:      signUpResponse.User.SignUpAt,
		},
	}
	assert.Equal(t, expected, signUpResponse)
}

func TestSignIn(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the SignIn method
	request := user.SignInRequest{
		Account:  "test_account0",
		Password: "test_password0",
	}
	_, err := userService.SignIn(request)
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
	expected := user.UserResponse{
		User: toUserDetails(mockUser0()),
	}
	assert.Equal(t, expected, userDetails)
}

func TestGetByAccount(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the GetByAccount method
	userResponse, err := userService.GetByAccount("test_account0")
	assert.Equal(t, nil, err)

	// Check the returned user details
	expected := user.UserResponse{
		User: toUserDetails(mockUser0()),
	}
	assert.Equal(t, expected, userResponse)
}

func TestSearchByAccount(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	// Call the List method
	request := user.UserSearchRequest{
		LikeAccount: "test_account",
		Limit:      10,
		Offset:     0,
		OrderBy:    "account",
	}
	UserSearchResponse, err := userService.SearchByAccount(request)
	assert.Equal(t, nil, err)

	// Check the returned users
	expected := user.UserSearchResponse{
		Total: 2,
		Users: []user.UserDetails{
			*toUserDetails(mockUser0()),
			*toUserDetails(mockUser1()),
		},
	}
	assert.Equal(t, expected, UserSearchResponse)
}

func TestUpdateUser(t *testing.T) {
	// Mock the user models
	mockUserModels := &MockUserModels{}

	// Create a new user service with the mock user models
	userService := user.NewUserServiceWithModel(mockUserModels)

	location := user.Point{Latitude: 39.9042, Longitude: 116.4074}

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
		Location:      &user.Point{Latitude: u.Location.Latitude, Longitude: u.Location.Longitude},
		Icon:          u.Icon,
		Role:          u.Role,
		LastLoginTime: u.LastLoginTime,
		SignUpAt:      u.CreatedAt,
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
		Content:   []byte("{\"type\":\"alarm\",\"resource\":{\"type\":\"image\",\"url\":\"test_url\"},\"device\":\"test_device0\"}"),
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
			Content: user.AlarmContent{
				Type: "alarm",
				Resource: user.Resource{
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
			Content: user.DeviceNotificationContent{
				Action: "ADD_DEVICE",
				Device: "test_device1",
			},
		},
	}
	assert.Equal(t, 2, len(messages))
	assert.Equal(t, expected[0], messages[0])
	assert.Equal(t, expected[1], messages[1])
}
