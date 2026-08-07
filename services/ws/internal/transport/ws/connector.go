package ws

import (
	"MyMessenger/services/ws/internal/service"
	"context"

	"github.com/coder/websocket"
)

type Connector struct {
	conn *websocket.Conn
}

func (c *Connector) toServiceType(mt websocket.MessageType) service.MessageType {
	switch mt {
	case websocket.MessageBinary:
		return service.BinaryType
	case websocket.MessageText:
		return service.TextType
	default:
		return service.TextType
	}
}

func (c *Connector) fromServiceType(st service.MessageType) websocket.MessageType {
	switch st {
	case service.BinaryType:
		return websocket.MessageBinary
	case service.TextType:
		return websocket.MessageText
	default:
		return websocket.MessageText
	}
}

func (c *Connector) Read(ctx context.Context) (service.MessageType, []byte, error) {
	wType, msg, err := c.conn.Read(ctx)
	if err != nil {
		status := websocket.CloseStatus(err)
		if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
			return service.TextType, nil, service.ErrConnNormalClosure
		}
		return service.TextType, nil, err
	}
	return c.toServiceType(wType), msg, nil
}

func (c *Connector) Write(ctx context.Context, msgType service.MessageType, msg []byte) error {
	return c.conn.Write(ctx, c.fromServiceType(msgType), msg)
}

func (c *Connector) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}
