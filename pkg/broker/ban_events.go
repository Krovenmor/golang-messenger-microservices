package broker

import (
	"time"

	"github.com/google/uuid"
)

const (
	BanRequest EventType = "ban_req"
	BanEvent   EventType = "ban_ev"
)

type BanReason string

const (
	TooManyRequests BanReason = "too_many_reqs"
)

type BanRequestDTO struct {
	Type    EventType         `json:"type"`
	Payload BanRequestPayload `json:"payload"`
}

type BanRequestPayload struct {
	UserId uuid.UUID `json:"userId"`
	Reason BanReason `json:"reason"`
}

type BanEventDTO struct {
	Type    EventType       `json:"type"`
	Payload BanEventPayload `json:"payload"`
}

type BanEventPayload struct {
	UserId uuid.UUID     `json:"userId"`
	Ttl    time.Duration `json:"ttl"`
}
