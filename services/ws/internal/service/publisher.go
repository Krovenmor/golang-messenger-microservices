package service

import (
	"MyMessenger/pkg/broker"
	"context"
)

type Publisher interface {
	PublishUserStatus(ctx context.Context, status broker.StatusEvent)
}
