package things

type Event struct {
	DeviceID string `json:"device_id"`
	Type     string `json:"type"`
	Data     []byte `json:"data"`
}
