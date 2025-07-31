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
	ProductID    string   `json:"product_id"`
	Version      string   `json:"version"`
	SDKVersion   string   `json:"sdk_version"`
	IP           string   `json:"ip"`
	State        state    `json:"state"`
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
