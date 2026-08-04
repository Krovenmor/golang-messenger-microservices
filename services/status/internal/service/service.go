package service

import (
	"MyMessenger/pkg/broker"
	"context"
)

type UserStatus struct {
	Status   broker.Status
	LastSeen int64
}

type StatusService interface {
	GetStatus(ctx context.Context, userId string) UserStatus
	SaveStatus(ctx context.Context, event broker.StatusEvent) error
}
