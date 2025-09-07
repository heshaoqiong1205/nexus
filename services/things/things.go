package things

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"nexus/models"
	"time"

	"github.com/google/uuid"
)

type ThingService struct {
	deviceModels       models.IDeviceModels
	licenseModels      models.ILicenseModels
	productModels      models.IProductModels
	desiredStateModels models.IDesiredStateModels
}

type ServiceRequest struct {
	DeviceID      string `json:"device_id"`
	Method        string `json:"method"`
	TransactionID string `json:"transaction_id"`
	Data          []byte `json:"data"`
}

type ServiceResponse struct {
	DeviceID      string `json:"device_id"`
	TransactionID string `json:"transaction_id"`
	Data          []byte `json:"data"`
}

type IResponse interface {
	SetResult(result bool)
	SetMessage(message string)
}

type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func (r *BaseResponse) SetResult(result bool) {
	r.Result = result
}

func (r *BaseResponse) SetMessage(message string) {
	r.Message = message
}

type ActiveRequest struct {
	LicenseID    string   `json:"license_id"`
	Authenticate string   `json:"authenticate"`
	ProductID    string   `json:"product_id"`
	Version      string   `json:"version"`
	SDKVersion   string   `json:"sdk_version"`
	IP           string   `json:"ip"`
	State        State    `json:"state"`
	Features     Features `json:"features"`
}

type DevicesQuery struct {
	GroupList []string `json:"group_list"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
	OrderBy   *string  `json:"order_by"`
}

type ConfirmDesiredStateRequest struct {
	LastID int64 `json:"last_id"`
}

func newServiceRequest(method string, deviceID string, data []byte) *ServiceRequest {
	return &ServiceRequest{
		DeviceID:      deviceID,
		Method:        method,
		TransactionID: uuid.NewString(),
		Data:          data,
	}
}

func NewThingService() *ThingService {
	return &ThingService{
		deviceModels:       &models.DeviceModels{},
		licenseModels:      &models.LicenseModels{},
		productModels:      &models.ProductModels{},
		desiredStateModels: &models.DesiredStateModels{},
	}
}

func NewThingServiceWithModels(deviceModels models.IDeviceModels,
	licenseModels models.ILicenseModels, productModels models.IProductModels,
	desiredStateModels models.IDesiredStateModels) *ThingService {
	return &ThingService{
		deviceModels:       deviceModels,
		licenseModels:      licenseModels,
		productModels:      productModels,
		desiredStateModels: desiredStateModels,
	}
}

func (service *ThingService) Active(request *ActiveRequest) (*IoTDevice, error) {

	err := service.authentication(request.LicenseID, request.Authenticate)
	if err != nil {
		return nil, err
	}
	product, err := service.getProduct(request.ProductID)
	if err != nil {
		return nil, err
	}
	var requiredFeatures []string
	err = json.Unmarshal(product.RequiredFeatures, &requiredFeatures)
	if err != nil {
		return nil, err
	}

	newFeature, err := service.checkFeatures(requiredFeatures, request.Features)
	if err != nil {
		return nil, err
	}

	bState, err := request.State.Marshal()
	if err != nil {
		return nil, err
	}

	device := &models.Device{
		ID:         models.GenerateDeviceID(),
		Name:       product.Name + "-" + uuid.NewString()[0:4],
		SecretKey:  models.GenerateSecretKey(),
		LicenseID:  request.LicenseID,
		ProductID:  request.ProductID,
		Features:   newFeature,
		Version:    request.Version,
		SDKVersion: request.SDKVersion,
		IP:         request.IP,
		State:      bState,
		UpdatedAt:  time.Now(),
		ActiveAt:   time.Now(),
		CreatedAt:  time.Now(),
		Status:     true,
	}
	err = service.deviceModels.Create(device)
	if err != nil {
		return nil, err
	}
	return NewIoTDevice(*device)
}

func (service *ThingService) Deactivate(deviceID string) error {
	if deviceID == "" {
		return errors.New("device id cannot be empty")
	}
	device, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return errors.New("device not found")
	}
	if !device.Status {
		return errors.New("device is not active")
	}
	device.Status = false
	device.UpdatedAt = time.Now()
	err = service.deviceModels.Update(&device)
	if err != nil {
		return errors.New("failed to update device status")
	}
	return nil
}

func (service *ThingService) GetDevice(deviceID string) (*IoTDevice, error) {
	if deviceID == "" {
		return nil, errors.New("device id cannot be empty")
	}
	device, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return nil, errors.New("device not found")
	}
	return NewIoTDevice(device)
}

func (service *ThingService) GetDevices(query DevicesQuery) ([]*IoTDevice, error) {
	devices, err := service.deviceModels.List(query.GroupList, query.Page, query.PageSize, query.OrderBy)
	if err != nil {
		return nil, err
	}
	var iotDevices []*IoTDevice
	for _, device := range devices {
		iotDevice, err := NewIoTDevice(device)
		if err != nil {

			return nil, err
		}
		iotDevices = append(iotDevices, iotDevice)
	}
	return iotDevices, nil
}

func (service *ThingService) RPC(deviceID, method string, payload interface{}, response IResponse) {
	response.SetResult(false)
	data, err := json.Marshal(payload)
	if err != nil {
		response.SetMessage(err.Error())
		return
	}
	req := newServiceRequest(method, deviceID, data)
	resp, err := service.rpc(req)
	if err != nil {
		response.SetMessage(err.Error())
		return
	}
	if response != nil {
		err = json.Unmarshal(resp.Data, response)
		if err != nil {
			response.SetMessage(err.Error())
			return
		}
		return
	}
}

func (service *ThingService) rpc(request *ServiceRequest) (*ServiceResponse, error) {
	device, err := service.deviceModels.Get(request.DeviceID)
	if err != nil {
		return nil, errors.New("device not found")
	}
	if !device.Status {
		return nil, errors.New("device is not active")
	}
	response := &ServiceResponse{
		DeviceID:      request.DeviceID,
		TransactionID: request.TransactionID,
		Data:          request.Data,
	}
	return response, nil
}

func (service *ThingService) HandleEvent(deviceID string, event Event) error {
	return nil
}

func (service *ThingService) HandleState(deviceID string, state State) error {
	// Validate input parameters
	if deviceID == "" {
		return errors.New("device ID cannot be empty")
	}

	// Validate incoming state
	if err := state.Validate(); err != nil {
		return errors.New("invalid state: " + err.Error())
	}

	// Get device from database
	device, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return errors.New("device not found: " + err.Error())
	}

	// Check if device is active
	if !device.Status {
		return errors.New("device is not active")
	}

	// Parse current device state
	var currentState State
	if len(device.State) > 0 {
		if err := json.Unmarshal(device.State, &currentState); err != nil {
			log.Printf("Warning: failed to parse current device state for %s: %v", deviceID, err)
			// Initialize with empty state if parsing fails
			currentState = State{}
		}
	}

	// Merge states - only update non-nil fields
	service.mergeStates(&currentState, state)

	// Marshal updated state
	updatedStateData, err := json.Marshal(&currentState)
	if err != nil {
		return errors.New("failed to marshal updated state: " + err.Error())
	}

	// Update device with new state and timestamp
	device.State = updatedStateData
	device.UpdatedAt = time.Now()

	// Save to database
	if err := service.deviceModels.Update(&device); err != nil {
		return errors.New("failed to update device state: " + err.Error())
	}

	log.Printf("Device %s state updated successfully", deviceID)
	return nil
}


func (service *ThingService) FetchDesiredState(deviceID string) (DesiredState, error) {
	// Validate device exists
	_, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return DesiredState{}, errors.New("device not found: " + err.Error())
	}

	// Get the desired state for this device (one per device)
	desiredState, err := service.desiredStateModels.Get(deviceID)
	if err != nil {
		// If no desired state exists, return empty state
		return DesiredState{}, nil
	}

	// Parse the desired state JSON into State struct
	var state DesiredState
	if err := json.Unmarshal(desiredState.State, &state); err != nil {
		return DesiredState{}, errors.New("failed to parse desired state: " + err.Error())
	}

	return state, nil
}

func (service *ThingService) UpdateDesiredState(deviceID string, state State) error {
	// Validate device exists
	device, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return errors.New("device not found: " + err.Error())
	}

	// Validate the state
	if err := state.Validate(); err != nil {
		return errors.New("invalid desired state: " + err.Error())
	}

	var originState State
	if err := json.Unmarshal(device.State, &originState); err != nil {
		return errors.New("failed to unmarshal original state: " + err.Error())
	}

	desiredState, err := service.desiredStateModels.Get(deviceID)
	if err != nil {
		// If no desired state exists, create a new one
		return  errors.New("desired state not found: " + err.Error())
	}

	// Create a desired state object with differences from the current state
	lastID, updateDesiredState := NewDesiredState(desiredState.LastDesiredID, &originState, &state)

	if lastID != desiredState.LastDesiredID {
		stateJSON, err := json.Marshal(updateDesiredState)
		if err != nil {
			return errors.New("failed to marshal desired state: " + err.Error())
		}
		// Update the desired state in the database
		desiredState.ID = deviceID
		desiredState.State = stateJSON
		desiredState.Status = true
		desiredState.LastDesiredID = lastID

		if err := service.desiredStateModels.Update(&desiredState); err != nil {
			return errors.New("failed to update desired state: " + err.Error())
		}
		log.Printf("Updated desired state for device %s", deviceID)

	}

	return nil
}

func (service *ThingService) ConfirmDesiredState(deviceID string, confirmRequest ConfirmDesiredStateRequest) error {
	for i := 0; i < 5; i++ {
		if err := service.confirmDesiredState(deviceID, confirmRequest); err != nil {
			log.Printf("Failed to confirm desired state for device %s (attempt %d): %v", deviceID, i+1, err)
			continue
		}
		return nil
	}
	return errors.New("failed to confirm desired state after 5 attempts")
}

// SetDesiredState sets a new desired state for a device
func (service *ThingService) SetDesiredState(deviceID string, state State) error {
	for i := 0; i < 5; i++ {
		if err := service.updateDesiredStateWithVersion(deviceID, state); err != nil {
			log.Printf("Failed to set desired state for device %s (attempt %d): %v", deviceID, i+1, err)
			continue
		}
		return nil
	}
	return errors.New("failed to set desired state after 5 attempts")
}

// GetDesiredStateHistory returns the history of desired states for a device
func (service *ThingService) GetDesiredStateHistory(deviceID string, page int, pageSize int) ([]models.DesiredState, error) {
	// Validate device exists
	_, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return nil, errors.New("device not found: " + err.Error())
	}

	return service.desiredStateModels.List([]string{deviceID}, page, pageSize, nil)
}

/*
internal function
*/
func (service *ThingService) authentication(licenseID string, authenticate string) error {
	license, err := service.licenseModels.Get(licenseID)
	if err != nil {
		return err
	}
	if !license.Status {
		return errors.New("license is not active")
	}
	h := hmac.New(sha256.New, []byte(license.Key))
	h.Write([]byte(license.ID))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	log.Printf("signature: %s, authenticate: %s\n", signature, authenticate)

	if signature == authenticate {
		return nil
	}
	return errors.New("invalid authentication")
}

func (service *ThingService) checkFeatures(requiredFeatures []string, features Features) ([]byte, error) {
	featuresJSON, err := features.Marshal()
	if err != nil {
		return nil, err
	}
	log.Printf("requiredFeatures: %+v, features: %+v\n", requiredFeatures, features)
	err = features.ValidateRequirements(requiredFeatures)
	if err != nil {
		return nil, err
	}
	return featuresJSON, nil
}

func (service *ThingService) getProduct(productID string) (models.Product, error) {
	if productID == "" {
		return models.Product{}, errors.New("product id cannot be empty")
	}
	product, err := service.productModels.Get(productID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

// mergeStates merges the incoming state into the current state
// Only non-nil fields from the incoming state will overwrite the current state
func (service *ThingService) mergeStates(current *State, incoming State) {
	if incoming.Cruise != nil {
		current.Cruise = incoming.Cruise
	}
	if incoming.DecibelDetection != nil {
		current.DecibelDetection = incoming.DecibelDetection
	}
	if incoming.MotionDetection != nil {
		current.MotionDetection = incoming.MotionDetection
	}
	if incoming.MotionTracking != nil {
		current.MotionTracking = incoming.MotionTracking
	}
	if incoming.NightVision != nil {
		current.NightVision = incoming.NightVision
	}
	if incoming.PrivacyMode != nil {
		current.PrivacyMode = incoming.PrivacyMode
	}
	if incoming.Record != nil {
		current.Record = incoming.Record
	}
	if incoming.Siren != nil {
		current.Siren = incoming.Siren
	}
	if incoming.Storage != nil {
		current.Storage = incoming.Storage
	}
	if incoming.Video != nil {
		current.Video = incoming.Video
	}
	if incoming.Volume != nil {
		current.Volume = incoming.Volume
	}
}

func (service *ThingService) confirmDesiredState(deviceID string, confirmRequest ConfirmDesiredStateRequest) error {
	// Validate device exists
	desiredState, err := service.desiredStateModels.Get(deviceID)
	if err != nil {
		return errors.New("device not found: " + err.Error())
	}

	var state DesiredState
	if err := json.Unmarshal(desiredState.State, &state); err != nil {
		return errors.New("failed to parse desired state: " + err.Error())
	}

	state.Confirm(confirmRequest.LastID)

	data, err := state.Marshal()
	if err != nil {
		return err
	}
	desiredState.State = data

	version := desiredState.Version
	desiredState.Version = desiredState.Version + 1
	// With simplified model, we just confirm the desired state for this device
	return service.desiredStateModels.UpdateWithVersion(&desiredState, version)
}

func (service *ThingService) updateDesiredStateWithVersion(deviceID string, state State) error {
	// Validate device exists
	device, err := service.deviceModels.Get(deviceID)
	if err != nil {
		return errors.New("device not found: " + err.Error())
	}

	if !device.Status {
		return errors.New("device is not active")
	}

	// Validate the state
	if err := state.Validate(); err != nil {
		return errors.New("invalid desired state: " + err.Error())
	}

	// Marshal state to JSON
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return errors.New("failed to marshal desired state: " + err.Error())
	}

	// Get current desired state to determine version
	currentVersion := 1
	if existingState, err := service.desiredStateModels.Get(deviceID); err == nil {
		currentVersion = existingState.Version + 1
	}

	// Create new desired state (or update existing one)
	newDesiredState := &models.DesiredState{
		ID:           deviceID, // DeviceID as primary key
		State: stateJSON,
		Version:      currentVersion,
		Status:       true,
	}

	// Try to update existing state first, create if doesn't exist
	err = service.desiredStateModels.Update(newDesiredState)
	if err != nil {
		// If update fails, try create
		if err := service.desiredStateModels.Create(newDesiredState); err != nil {
			return errors.New("failed to create desired state: " + err.Error())
		}
	}

	log.Printf("Set desired state for device %s (version %d)", deviceID, currentVersion)
	return nil
}
