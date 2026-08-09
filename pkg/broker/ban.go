package broker

import "github.com/google/uuid"

const (
	BanEvent EventType = "ban"
)

type BanReason string

const (
	TooManyRequests BanReason = "too_many_reqs"
)

type BanEventDTO struct {
	Type    EventType       `json:"type"`
	Payload BanEventPayload `json:"payload"`
}

type BanEventPayload struct {
	UserId uuid.UUID `json:"userId"`
	Reason BanReason `json:"reason"`
}
