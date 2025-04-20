package things

type ServiceRequest struct {
	DeviceID      string `json:"device_id"`
	TransactionID string `json:"transaction_id"`
	Data          []byte `json:"data"`
}

type ServiceResponse struct {
	DeviceID      string `json:"device_id"`
	TransactionID string `json:"transaction_id"`
	Result        string `json:"result"`
	Data          []byte `json:"data"`
}
