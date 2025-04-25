package services_testing

import (
	"encoding/json"
	"nexus/models"
	"nexus/services/things"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDeviceModels struct {
	mock.Mock
}

func (m *MockDeviceModels) Get(licenseID string) (models.Device, error) {
	return models.Device{}, nil
}

func (m *MockDeviceModels) List(group_list []string, page int, pageSize int, orderBY string) ([]models.Device, error) {
	return []models.Device{}, nil
}

func (m *MockDeviceModels) Create(device *models.Device) error {
	return nil
}

func (m *MockDeviceModels) Update(device *models.Device) error {
	return nil
}

func (m *MockDeviceModels) Delete(id string) error {
	return nil
}

type MockLicenseModels struct {
	mock.Mock
}

func (m *MockLicenseModels) Get(licenseID string) (models.License, error) {
	result := m.Called(licenseID)
	if result.Get(0) == nil {
		return models.License{}, result.Error(1)
	}
	return result.Get(0).(models.License), nil
}

func (m *MockLicenseModels) Create(license *models.License) error {
	return nil
}

func (m *MockLicenseModels) Update(license *models.License) error {
	return nil
}

type MockProductModels struct {
	mock.Mock
}

func (m *MockProductModels) Get(id string) (models.Product, error) {
	result := m.Called(id)
	if result.Get(0) == nil {
		return models.Product{}, result.Error(1)
	}
	return result.Get(0).(models.Product), nil
}

func (m *MockProductModels) Create(product *models.Product) error {
	return nil
}

func (m *MockProductModels) Update(product *models.Product) error {
	return nil
}

func (m *MockProductModels) Delete(id string) error {
	return nil
}

func TestAcitve(t *testing.T) {
	requestFeatures, _ := json.Marshal([]string{"p2p", "webrtc", "upnp", "ai"})
	mockLicenseModels := new(MockLicenseModels)
	mockDeviceModels := new(MockDeviceModels)
	mockProductModels := new(MockProductModels)
	mockLicenseModels.On("Get", "1234567890").Return(models.License{
		ID:        "1234567890",
		Key:       "test_key",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    true,
	}).Times(1)
	mockProductModels.On("Get", "234512345").Return(models.Product{
		ID:               "234512345",
		Name:             "test_product",
		Description:      "test_description",
		RequiredFeatures: requestFeatures,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           true,
	}).Times(1)

	service := things.NewThingServiceWithModels(mockDeviceModels, mockLicenseModels, mockProductModels)
	standard := "standard"
	upnp := "IDGV1"
	ai := "local"
	request := &things.ActiveRequest{
		LicenseID:    "1234567890",
		Authenticate: "plt41RL8aFt4SeTPSvDOVwKHFWPuJ1mVMBTr6wkERx8=",
		ProductID:    "234512345",
		Features: things.Features{
			P2P:    &standard,
			WebRTC: []string{"SRTP", "DC"},
			UPNP:   &upnp,
			AI:     &ai,
			VideoFeature: &things.VideoFeature{
				Num:              1,
				Codecs:           []string{"h264", "h265"},
				ResolutionRatios: []things.ResolutionRatio{{Width: 1920, Height: 1080}},
				Streams:          []int{1, 2},
			},
			AudioFeature: &things.AudioFeature{
				Codecs:       []string{"PCMA", "PCMU"},
				SamplingRate: 16000,
				Channels:     2,
				Bits:         16,
			},
		},
	}
	response, err := service.Active(request)
	assert.Equal(t, nil, err, "err should be nil")
	t.Logf("%v ", response)
}
