package ws

import (
	"MyMessenger/services/ws/internal/config"
	"MyMessenger/services/ws/internal/service"
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const (
	ErrWebSocTooManyReq = 4001
)

type WSHandler struct {
	wsService service.WsService
	readLimit int64
}

func NewWSHandler(wsService service.WsService, conf *config.WsConfig) *WSHandler {
	return &WSHandler{
		wsService: wsService,
		readLimit: conf.ReadLimit,
	}
}

func (h *WSHandler) HandleConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID, aToken string) {
	ctx := r.Context()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("HandleConnection: Trouble with Accept: %q", err.Error())
		return
	}
	conn.SetReadLimit(h.readLimit)

	connCloseCode := websocket.StatusNormalClosure
	connCloseMsg := "Normal closure"

	err = h.wsService.StartService(ctx, &Connector{conn: conn}, userId, aToken)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			connCloseCode = websocket.StatusPolicyViolation
			connCloseMsg = "unauthorized"

		case errors.Is(err, service.ErrTooManyRequests):
			connCloseCode = ErrWebSocTooManyReq
			connCloseMsg = "too many requests"

		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			connCloseCode = websocket.StatusGoingAway
			connCloseMsg = "server shutdown or timeout"

		default:
			connCloseCode = websocket.StatusInternalError
			connCloseMsg = "internal error"
		}
	}

	err = conn.Close(connCloseCode, connCloseMsg)
	if err != nil {
		status := websocket.CloseStatus(err)
		if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
			return
		}
		if errors.Is(err, net.ErrClosed) {
			return
		}
		log.Printf("conn.Close() err for user %s: %v", userId, err)
	}
}
