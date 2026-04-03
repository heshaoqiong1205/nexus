package things

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"nexus/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	existingModel, existingDesired, hasExisting, err := service.loadDesiredState(deviceID)
	if err != nil {
		return err
	}

	if hasExisting {
		nextDesired, changed := reconcileDesiredWithReported(existingDesired, &currentState)
		if changed {
			stateJSON, err := nextDesired.Marshal()
			if err != nil {
				return errors.New("failed to marshal desired state: " + err.Error())
			}
			existingModel.State = stateJSON
			version := existingModel.Version
			existingModel.Version++
			if err := service.desiredStateModels.UpdateWithVersion(&existingModel, version); err != nil {
				return errors.New("failed to update desired state: " + err.Error())
			}
		}
	}

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no desired state exists, return empty state
			return DesiredState{}, nil
		}
		return DesiredState{}, errors.New("failed to get desired state: " + err.Error())
	}

	// Parse the desired state JSON into State struct
	var state DesiredState
	if err := json.Unmarshal(desiredState.State, &state); err != nil {
		return DesiredState{}, errors.New("failed to parse desired state: " + err.Error())
	}

	return state, nil
}

func (service *ThingService) UpdateDesiredState(deviceID string, state State) error {
	return service.updateDesiredStateWithVersion(deviceID, state)
}

// SetDesiredState sets a new desired state for a device
func (service *ThingService) SetDesiredState(deviceID string, state State) error {
	var lastErr error
	for i := 0; i < 5; i++ {
		if err := service.updateDesiredStateWithVersion(deviceID, state); err != nil {
			lastErr = err
			if !isDesiredStateVersionConflict(err) {
				return err
			}
			log.Printf("Failed to set desired state for device %s (attempt %d): %v", deviceID, i+1, err)
			continue
		}
		return nil
	}
	if lastErr != nil {
		return errors.New("failed to set desired state after 5 attempts")
	}
	return nil
}

// GetDesiredState returns the desired state records for a device
func (service *ThingService) GetDesiredState(deviceID string, page int, pageSize int) ([]models.DesiredState, error) {
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

	existingModel, existingDesired, hasExisting, err := service.loadDesiredState(deviceID)
	if err != nil {
		return err
	}

	lastID := int64(0)
	lastEpoch := int64(0)
	if hasExisting {
		lastID = existingModel.LastDesiredID
		lastEpoch = existingDesired.MaxEpoch()
	}

	nextDesired, nextLastID, _, changed := mergeDesiredStateUpdate(existingDesired, lastEpoch, lastID, &state)
	if !changed {
		return nil
	}

	stateJSON, err := nextDesired.Marshal()
	if err != nil {
		return errors.New("failed to marshal desired state: " + err.Error())
	}

	if hasExisting {
		existingModel.State = stateJSON
		existingModel.Status = true
		existingModel.LastDesiredID = nextLastID
		version := existingModel.Version
		existingModel.Version++
		if err := service.desiredStateModels.UpdateWithVersion(&existingModel, version); err != nil {
			return errors.New("failed to update desired state: " + err.Error())
		}
		log.Printf("Updated desired state for device %s (version %d)", deviceID, existingModel.Version)
		return nil
	}

	newDesiredState := &models.DesiredState{
		ID:            deviceID,
		State:         stateJSON,
		Version:       1,
		Status:        true,
		LastDesiredID: nextLastID,
	}
	if err := service.desiredStateModels.Create(newDesiredState); err != nil {
		return errors.New("failed to create desired state: " + err.Error())
	}

	log.Printf("Created desired state for device %s", deviceID)
	return nil
}

func mergeDesiredStateUpdate(existing DesiredState, startEpoch, startID int64, update *State) (DesiredState, int64, int64, bool) {
	next := existing
	changed := false
	currentEpoch := startEpoch
	currentID := startID

	if update.Video != nil {
		if !desiredVideoMatches(next.Video, update.Video) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Video: update.Video})
			next.Video = diff.Video
			changed = true
		}
	}

	if update.Storage != nil {
		if !desiredStorageMatches(next.Storage, update.Storage) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Storage: update.Storage})
			next.Storage = diff.Storage
			changed = true
		}
	}

	if update.Record != nil {
		if !desiredRecordMatches(next.Record, update.Record) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Record: update.Record})
			next.Record = diff.Record
			changed = true
		}
	}

	if update.MotionDetection != nil {
		if !desiredVMDMatches(next.MotionDetection, update.MotionDetection) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{MotionDetection: update.MotionDetection})
			next.MotionDetection = diff.MotionDetection
			changed = true
		}
	}

	if update.DecibelDetection != nil {
		if !desiredDetectionMatches(next.DecibelDetection, update.DecibelDetection) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{DecibelDetection: update.DecibelDetection})
			next.DecibelDetection = diff.DecibelDetection
			changed = true
		}
	}

	if update.Cruise != nil {
		if !desiredCruiseMatches(next.Cruise, update.Cruise) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Cruise: update.Cruise})
			next.Cruise = diff.Cruise
			changed = true
		}
	}

	if update.Siren != nil {
		if !desiredSirenMatches(next.Siren, update.Siren) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Siren: update.Siren})
			next.Siren = diff.Siren
			changed = true
		}
	}

	if update.Volume != nil {
		if !desiredVolumeMatches(next.Volume, update.Volume) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{Volume: update.Volume})
			next.Volume = diff.Volume
			changed = true
		}
	}

	if update.PrivacyMode != nil {
		if !desiredBooleanMatches(next.PrivacyMode, update.PrivacyMode) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{PrivacyMode: update.PrivacyMode})
			next.PrivacyMode = diff.PrivacyMode
			changed = true
		}
	}

	if update.NightVision != nil {
		if !desiredBooleanMatches(next.NightVision, update.NightVision) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{NightVision: update.NightVision})
			next.NightVision = diff.NightVision
			changed = true
		}
	}

	if update.MotionTracking != nil {
		if !desiredBooleanMatches(next.MotionTracking, update.MotionTracking) {
			var diff DesiredState
			currentID, currentEpoch, diff = NewDesiredStateWithEpoch(currentEpoch, currentID, &State{MotionTracking: update.MotionTracking})
			next.MotionTracking = diff.MotionTracking
			changed = true
		}
	}

	return next, currentID, currentEpoch, changed
}

func (service *ThingService) loadDesiredState(deviceID string) (models.DesiredState, DesiredState, bool, error) {
	existingModel, err := service.desiredStateModels.Get(deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.DesiredState{}, DesiredState{}, false, nil
		}
		return models.DesiredState{}, DesiredState{}, false, errors.New("failed to load desired state: " + err.Error())
	}

	var existingDesired DesiredState
	if len(existingModel.State) > 0 {
		if err := json.Unmarshal(existingModel.State, &existingDesired); err != nil {
			return models.DesiredState{}, DesiredState{}, false, errors.New("failed to parse existing desired state: " + err.Error())
		}
	}

	return existingModel, existingDesired, true, nil
}

func desiredVideoMatches(desired *DesiredVideo, update *Video) bool {
	return desired != nil && desired.Video.Equal(update)
}

func desiredStorageMatches(desired *DesiredStorage, update *Storage) bool {
	return desired != nil && desired.Storage.Equal(update)
}

func desiredRecordMatches(desired *DesiredRecord, update *Record) bool {
	return desired != nil && desired.Record.Equal(update)
}

func desiredVMDMatches(desired *DesiredVMD, update *VMD) bool {
	return desired != nil && desired.VMD.Equal(update)
}

func desiredDetectionMatches(desired *DesiredDecibelDetection, update *Detection) bool {
	return desired != nil && desired.Detection.Equal(update)
}

func desiredCruiseMatches(desired *DesiredCruise, update *Cruise) bool {
	return desired != nil && desired.Cruise.Equal(update)
}

func desiredSirenMatches(desired *DesiredSiren, update *Siren) bool {
	return desired != nil && desired.Siren.Equal(update)
}

func desiredVolumeMatches(desired *DesiredVolume, update *IntValue) bool {
	return desired != nil && intPtrEqual(desired.Value, intValue(update))
}

func desiredBooleanMatches(desired *DesiredBoolean, update *BoolValue) bool {
	return desired != nil && boolPtrEqual(desired.Value, boolValue(update))
}

func reportedHasVersion(meta StateMeta) bool {
	return meta.ID != nil && meta.Epoch != nil
}

func reconcileDesiredWithReported(existing DesiredState, reported *State) (DesiredState, bool) {
	next := existing
	changed := false

	if reported.Video != nil && next.Video != nil && reportedHasVersion(reported.Video.StateMeta) &&
		shouldDeleteDesired(*reported.Video.Epoch, *reported.Video.ID, next.Video.Epoch, next.Video.ID) {
		next.Video = nil
		changed = true
	}

	if reported.Storage != nil && next.Storage != nil && reportedHasVersion(reported.Storage.StateMeta) &&
		shouldDeleteDesired(*reported.Storage.Epoch, *reported.Storage.ID, next.Storage.Epoch, next.Storage.ID) {
		next.Storage = nil
		changed = true
	}

	if reported.Record != nil && next.Record != nil && reportedHasVersion(reported.Record.StateMeta) &&
		shouldDeleteDesired(*reported.Record.Epoch, *reported.Record.ID, next.Record.Epoch, next.Record.ID) {
		next.Record = nil
		changed = true
	}

	if reported.MotionDetection != nil && next.MotionDetection != nil && reportedHasVersion(reported.MotionDetection.StateMeta) &&
		shouldDeleteDesired(*reported.MotionDetection.Epoch, *reported.MotionDetection.ID, next.MotionDetection.Epoch, next.MotionDetection.ID) {
		next.MotionDetection = nil
		changed = true
	}

	if reported.DecibelDetection != nil && next.DecibelDetection != nil && reportedHasVersion(reported.DecibelDetection.StateMeta) &&
		shouldDeleteDesired(*reported.DecibelDetection.Epoch, *reported.DecibelDetection.ID, next.DecibelDetection.Epoch, next.DecibelDetection.ID) {
		next.DecibelDetection = nil
		changed = true
	}

	if reported.Cruise != nil && next.Cruise != nil && reportedHasVersion(reported.Cruise.StateMeta) &&
		shouldDeleteDesired(*reported.Cruise.Epoch, *reported.Cruise.ID, next.Cruise.Epoch, next.Cruise.ID) {
		next.Cruise = nil
		changed = true
	}

	if reported.Siren != nil && next.Siren != nil && reportedHasVersion(reported.Siren.StateMeta) &&
		shouldDeleteDesired(*reported.Siren.Epoch, *reported.Siren.ID, next.Siren.Epoch, next.Siren.ID) {
		next.Siren = nil
		changed = true
	}

	if reported.Volume != nil && next.Volume != nil && reportedHasVersion(reported.Volume.StateMeta) &&
		shouldDeleteDesired(*reported.Volume.Epoch, *reported.Volume.ID, next.Volume.Epoch, next.Volume.ID) {
		next.Volume = nil
		changed = true
	}

	if reported.PrivacyMode != nil && next.PrivacyMode != nil && reportedHasVersion(reported.PrivacyMode.StateMeta) &&
		shouldDeleteDesired(*reported.PrivacyMode.Epoch, *reported.PrivacyMode.ID, next.PrivacyMode.Epoch, next.PrivacyMode.ID) {
		next.PrivacyMode = nil
		changed = true
	}

	if reported.NightVision != nil && next.NightVision != nil && reportedHasVersion(reported.NightVision.StateMeta) &&
		shouldDeleteDesired(*reported.NightVision.Epoch, *reported.NightVision.ID, next.NightVision.Epoch, next.NightVision.ID) {
		next.NightVision = nil
		changed = true
	}

	if reported.MotionTracking != nil && next.MotionTracking != nil && reportedHasVersion(reported.MotionTracking.StateMeta) &&
		shouldDeleteDesired(*reported.MotionTracking.Epoch, *reported.MotionTracking.ID, next.MotionTracking.Epoch, next.MotionTracking.ID) {
		next.MotionTracking = nil
		changed = true
	}

	return next, changed
}

func isDesiredStateVersionConflict(err error) bool {
	return err != nil && strings.Contains(err.Error(), "desired state version conflict")
}
