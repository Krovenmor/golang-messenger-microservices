package service

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Requests int

const (
	MinReq Requests = iota

	NewStatusReq
	SubUserReq

	MaxReq
)

type Request struct {
	Req     Requests        `json:"req"`
	Payload json.RawMessage `json:"payload"`
}

type RPNewStatus struct {
	NewStatus int `json:"newStatus"`
}

type RPSubUser struct {
	UserId uuid.UUID `json:"userId"`
}
