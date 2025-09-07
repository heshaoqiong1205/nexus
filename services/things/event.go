package things

type Event struct {
	Timestamp int64`json:"timestamp"`
	Type      string    `json:"type"`
	Data      []byte    `json:"data"`
}
