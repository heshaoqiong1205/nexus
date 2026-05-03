package mqtt

import "encoding/json"

type StatusMessage struct {
	Online    bool   `json:"online"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type RPCRequestEnvelope struct {
	ID        string          `json:"id"`
	DeviceID  string          `json:"device_id"`
	Method    string          `json:"method"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
	TimeoutMS int64           `json:"timeout_ms,omitempty"`
}

type RPCResponseEnvelope struct {
	ID        string          `json:"id"`
	DeviceID  string          `json:"device_id,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     string          `json:"error,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
}
