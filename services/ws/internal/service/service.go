package service

import (
	"context"

	"github.com/google/uuid"
)

type WsService interface {
	StartService(ctx context.Context, conn Connector, userId uuid.UUID, accessToken string) error
}
