package service

import (
	"MyMessenger/pkg/broker"
	"context"
	"time"
)

type StatusRepo interface {
	GetStatus(ctx context.Context, userId string) (*UserStatus, error)
	SaveStatus(ctx context.Context, status broker.StatusPayload, ttl time.Duration) error
}
