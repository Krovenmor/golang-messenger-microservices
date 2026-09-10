package broker

import "encoding/json"

type EventType string

type Event struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type EventRaw struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
