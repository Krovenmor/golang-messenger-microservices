package service

import (
	"context"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WsService interface {
	StartService(ctx context.Context, conn *websocket.Conn, userId uuid.UUID)
}
