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
	deviceModels  models.IDeviceModels
	licenseModels models.ILicenseModels
	productModels models.IProductModels
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
	DeviceID     string   `json:"device_id"`
	ProductID    string   `json:"product_id"`
	Version      string   `json:"version"`
	SDKVersion   string   `json:"sdk_version"`
	IP           string   `json:"ip"`
	State        state    `json:"state"`
	Features     Features `json:"featrues"`
}

type DevicesQuery struct {
	GroupList []string `json:"group_id"`
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
	OrderBy   string   `json:"order_by"`
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
		deviceModels:  &models.DeviceModels{},
		licenseModels: &models.LicenseModels{},
		productModels: &models.ProductModels{},
	}
}

func NewThingServiceWithModels(deviceModels models.IDeviceModels,
	licenseModels models.ILicenseModels, productModels models.IProductModels) *ThingService {
	return &ThingService{
		deviceModels,
		licenseModels,
		productModels,
	}
}

func (service *ThingService) Active(request *ActiveRequest) (*IoTDevice, error) {

	err := service.authentication(request.LicenseID, request.Authenticate)
	if err != nil {
		return nil, err
	}
	requiredFeatures, err := service.getRequiredFeatures(request.ProductID)
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
	if request.DeviceID != "" {
		device, err := service.deviceModels.Get(request.DeviceID)
		if err != nil {
			return nil, errors.New("device not found")
		}
		if device.LicenseID != request.LicenseID {
			return nil, errors.New("device not belong to this license")
		}
		device.ProductID = request.ProductID
		device.Features = newFeature
		device.Version = request.Version
		device.SDKVersion = request.SDKVersion
		device.IP = request.IP
		device.State = bState
		device.UpdatedAt = time.Now()
		device.ActiveAt = time.Now()
		err = service.deviceModels.Update(&device)
		if err != nil {
			return nil, err
		}
		return NewIoTDevice(device)
	} else {
		device := &models.Device{
			ID:         request.DeviceID,
			ProductID:  request.ProductID,
			GroupID:    request.LicenseID,
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
		log.Println("========???========")
		return NewIoTDevice(*device)
	}
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

func (service *ThingService) HandleEvent(event Event) error {
	return nil
}

func (service *ThingService) UpdateState(state state) error {
	return nil
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
	if signature == authenticate {
		return nil
	}
	return errors.New("invalid token")
}

func (service *ThingService) checkFeatures(requiredFeatures []string, features Features) ([]byte, error) {
	featuresJSON, err := features.Marshal()
	if err != nil {
		return nil, err
	}
	err = features.ValidateRequirements(requiredFeatures)
	if err != nil {
		return nil, err
	}
	return featuresJSON, nil
}

func (service *ThingService) getRequiredFeatures(productID string) ([]string, error) {
	product, err := service.productModels.Get(productID)
	if err != nil {
		return nil, err
	}
	var requiredFeatures []string
	log.Printf("product.RequiredFeatures %v", product.RequiredFeatures)
	err = json.Unmarshal(product.RequiredFeatures, &requiredFeatures)
	if err != nil {
		return nil, err
	}
	return requiredFeatures, nil
}
