package things

import "encoding/json"

type Event struct {
	Timestamp int64           `json:"timestamp"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
}
