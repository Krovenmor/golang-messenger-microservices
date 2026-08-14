package broker

type EventType string

type Event struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}
