package services_testing

import (
	"encoding/json"
	"nexus/models"
	"nexus/services/things"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func boolPtr(v bool) *bool { return &v }
func intPtr(v int) *int    { return &v }

func reportedVideo(flip, osd bool, brightness, sharpness int) *things.Video {
	return &things.Video{
		Flip:       boolPtr(flip),
		OSD:        boolPtr(osd),
		Brightness: intPtr(brightness),
		Sharpness:  intPtr(sharpness),
	}
}

func reportedInt(value int) *things.IntValue {
	return &things.IntValue{Value: intPtr(value)}
}

func reportedBool(value bool) *things.BoolValue {
	return &things.BoolValue{Value: boolPtr(value)}
}

type MockDeviceModels struct {
	mock.Mock
}

func (m *MockDeviceModels) Get(deviceID string) (models.Device, error) {
	args := m.Called(deviceID)
	if args.Get(0) == nil {
		return models.Device{}, args.Error(1)
	}
	return args.Get(0).(models.Device), args.Error(1)
}

func (m *MockDeviceModels) List(group_list []string, page int, pageSize int, orderBY *string) ([]models.Device, error) {
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

type MockDesiredStateModels struct {
	mock.Mock
}

func (m *MockDesiredStateModels) Create(desiredState *models.DesiredState) error {
	args := m.Called(desiredState)
	return args.Error(0)
}

func (m *MockDesiredStateModels) Get(deviceID string) (models.DesiredState, error) {
	args := m.Called(deviceID)
	if args.Get(0) == nil {
		return models.DesiredState{}, args.Error(1)
	}
	return args.Get(0).(models.DesiredState), args.Error(1)
}

func (m *MockDesiredStateModels) Update(desiredState *models.DesiredState) error {
	args := m.Called(desiredState)
	return args.Error(0)
}

func (m *MockDesiredStateModels) UpdateWithVersion(desiredState *models.DesiredState, version int) error {
	args := m.Called(desiredState, version)
	return args.Error(0)
}

func (m *MockDesiredStateModels) Delete(deviceID string) error {
	args := m.Called(deviceID)
	return args.Error(0)
}

func (m *MockDesiredStateModels) List(deviceIDs []string, page int, pageSize int, orderBy *string) ([]models.DesiredState, error) {
	args := m.Called(deviceIDs, page, pageSize, orderBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.DesiredState), args.Error(1)
}

func (m *MockDesiredStateModels) Count(deviceIDs []string) (int64, error) {
	args := m.Called(deviceIDs)
	return args.Get(0).(int64), args.Error(1)
}

func TestActive(t *testing.T) {
	requestFeatures, _ := json.Marshal([]string{"p2p", "webrtc", "upnp", "ai"})
	mockLicenseModels := new(MockLicenseModels)
	mockDeviceModels := new(MockDeviceModels)
	mockProductModels := new(MockProductModels)
	mockDesiredStateModels := new(MockDesiredStateModels)
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

	service := things.NewThingServiceWithModels(mockDeviceModels, mockLicenseModels, mockProductModels, mockDesiredStateModels)
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
				Codecs:           []string{"H264", "H265"},
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

func TestFetchDesiredStateDeviceNotFound(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup mock to return device not found error
	mockDeviceModels.On("Get", "non-existent-device").Return(nil, assert.AnError)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	_, err := service.FetchDesiredState("non-existent-device")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not found")

	mockDeviceModels.AssertExpectations(t)
}

func TestFetchDesiredStateNoDesiredState(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup device exists
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: true,
	}, nil)

	// Setup no desired state exists
	mockDesiredStateModels.On("Get", "device123").Return(nil, gorm.ErrRecordNotFound)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	desiredState, err := service.FetchDesiredState("device123")
	assert.NoError(t, err)
	assert.Equal(t, things.DesiredState{}, desiredState)

	mockDeviceModels.AssertExpectations(t)
	mockDesiredStateModels.AssertExpectations(t)
}

func TestFetchDesiredStateSuccess(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Create sample desired state JSON
	desiredStateJSON := `{
		"video": {
			"epoch": 0,
			"id": 1,
			"flip": false,
			"osd": true,
			"brightness": 75,
			"sharpness": 50
		},
		"privacy_mode": {"epoch": 0, "id": 2, "value": false}
	}`

	// Setup device exists
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: true,
	}, nil)

	// Setup desired state exists
	mockDesiredStateModels.On("Get", "device123").Return(models.DesiredState{
		ID:            "device123",
		State:         []byte(desiredStateJSON),
		LastDesiredID: 2,
		Version:       1,
		Status:        true,
	}, nil)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	desiredState, err := service.FetchDesiredState("device123")
	assert.NoError(t, err)
	assert.NotNil(t, desiredState.Video)
	assert.Equal(t, int64(1), desiredState.Video.ID)
	assert.NotNil(t, desiredState.PrivacyMode)
	assert.Equal(t, int64(2), desiredState.PrivacyMode.ID)
	assert.Equal(t, int64(0), desiredState.PrivacyMode.Epoch)

	mockDeviceModels.AssertExpectations(t)
	mockDesiredStateModels.AssertExpectations(t)
}

func TestSetDesiredStateDeviceNotFound(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup device not found
	mockDeviceModels.On("Get", "non-existent-device").Return(nil, assert.AnError).Once()

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	state := things.State{Video: reportedVideo(false, true, 75, 50)}

	err := service.SetDesiredState("non-existent-device", state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not found")

	mockDeviceModels.AssertExpectations(t)
}

func TestSetDesiredStateDeviceInactive(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup inactive device
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: false, // Inactive device
	}, nil).Once()

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	state := things.State{Video: reportedVideo(false, true, 75, 50)}

	err := service.SetDesiredState("device123", state)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device is not active")

	mockDeviceModels.AssertExpectations(t)
}

func TestSetDesiredStateCreateNew(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup active device
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: true,
	}, nil)

	// Setup no existing desired state
	mockDesiredStateModels.On("Get", "device123").Return(nil, gorm.ErrRecordNotFound)

	mockDesiredStateModels.On("Create", mock.MatchedBy(func(desiredState *models.DesiredState) bool {
		return desiredState.ID == "device123" &&
			desiredState.Version == 1 &&
			desiredState.LastDesiredID == 1 &&
			desiredState.Status
	})).Return(nil)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	state := things.State{Video: reportedVideo(false, true, 75, 50)}

	err := service.SetDesiredState("device123", state)
	assert.NoError(t, err)

	mockDeviceModels.AssertExpectations(t)
	mockDesiredStateModels.AssertExpectations(t)
}

func TestSetDesiredStateUpdateExisting(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup active device
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: true,
	}, nil)

	// Setup existing desired state
	mockDesiredStateModels.On("Get", "device123").Return(models.DesiredState{
		ID:            "device123",
		Version:       1,
		LastDesiredID: 4,
		State:         []byte(`{"privacy_mode":{"epoch":0,"id":4,"value":false}}`),
	}, nil)

	mockDesiredStateModels.On("UpdateWithVersion", mock.MatchedBy(func(desiredState *models.DesiredState) bool {
		return desiredState.ID == "device123" &&
			desiredState.Version == 2 &&
			desiredState.LastDesiredID == 5 &&
			desiredState.Status
	}), 1).Return(nil)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	state := things.State{Video: reportedVideo(false, true, 75, 50)}

	err := service.SetDesiredState("device123", state)
	assert.NoError(t, err)

	mockDeviceModels.AssertExpectations(t)
	mockDesiredStateModels.AssertExpectations(t)
}

func TestNewDesiredStateCreatesFieldDiffs(t *testing.T) {
	originPrivacyMode := true
	updatePrivacyMode := false
	originVolume := 40
	updateVolume := 80
	flip := false
	osd := true
	brightness := 75
	sharpness := 50

	lastID, desiredState := things.NewDesiredState(7, &things.State{
		Volume:      reportedInt(originVolume),
		PrivacyMode: reportedBool(originPrivacyMode),
	}, &things.State{
		Video: reportedVideo(flip, osd, brightness, sharpness),
		Volume:      reportedInt(updateVolume),
		PrivacyMode: reportedBool(updatePrivacyMode),
	})

	assert.Equal(t, int64(10), lastID)
	assert.NotNil(t, desiredState.Video)
	assert.Equal(t, int64(8), desiredState.Video.ID)
	assert.Equal(t, int64(0), desiredState.Video.Epoch)
	assert.NotNil(t, desiredState.Volume)
	assert.Equal(t, int64(9), desiredState.Volume.ID)
	assert.NotNil(t, desiredState.PrivacyMode)
	assert.Equal(t, int64(10), desiredState.PrivacyMode.ID)
}

func TestGetDesiredStateHistoryDeviceNotFound(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup device not found
	mockDeviceModels.On("Get", "non-existent-device").Return(nil, assert.AnError)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	_, err := service.GetDesiredState("non-existent-device", 1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device not found")

	mockDeviceModels.AssertExpectations(t)
}

func TestGetDesiredStateHistorySuccess(t *testing.T) {
	mockDeviceModels := new(MockDeviceModels)
	mockDesiredStateModels := new(MockDesiredStateModels)

	// Setup device exists
	mockDeviceModels.On("Get", "device123").Return(models.Device{
		ID:     "device123",
		Status: true,
	}, nil)

	// Setup desired state history
	expectedHistory := []models.DesiredState{
		{
			ID:            "device123",
			State:         []byte(`{"video": {"id": 1, "bitrate": 1000}}`),
			LastDesiredID: 1,
			Version:       1,
			Status:        true,
			CreatedAt:     time.Now().Add(-time.Hour),
		},
		{
			ID:            "device123",
			State:         []byte(`{"privacy_mode": {"id": 2, "value": false}}`),
			LastDesiredID: 2,
			Version:       2,
			Status:        true,
			CreatedAt:     time.Now(),
		},
	}

	mockDesiredStateModels.On("List", []string{"device123"}, 1, 10, mock.AnythingOfType("*string")).Return(expectedHistory, nil)

	service := things.NewThingServiceWithModels(
		mockDeviceModels,
		new(MockLicenseModels),
		new(MockProductModels),
		mockDesiredStateModels,
	)

	history, err := service.GetDesiredState("device123", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, "device123", history[0].ID)
	assert.Equal(t, "device123", history[1].ID)

	mockDeviceModels.AssertExpectations(t)
	mockDesiredStateModels.AssertExpectations(t)
}
