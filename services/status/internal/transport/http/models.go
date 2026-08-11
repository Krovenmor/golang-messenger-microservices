package http

import "MyMessenger/services/status/internal/service"

type StatusResponseBody struct {
	Status   int   `json:"status"`
	LastSeen int64 `json:"lastSeen"`
}

func FromServiceStatus(us *service.UserStatus) *StatusResponseBody {
	return &StatusResponseBody{
		Status:   int(us.Status),
		LastSeen: us.LastSeen,
	}
}
