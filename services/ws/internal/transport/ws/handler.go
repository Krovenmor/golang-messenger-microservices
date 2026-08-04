package ws

import (
	"MyMessenger/services/ws/internal/config"
	"MyMessenger/services/ws/internal/service"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WSHandler struct {
	wsService service.WsService
	readLimit int64
}

func NewWSHandler(wsService service.WsService, conf *config.WsConfig) *WSHandler {
	return &WSHandler{
		wsService: wsService,
		readLimit: conf.WsReadLimit,
	}
}

func (h *WSHandler) HandleConnection(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	ctx := r.Context()

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("HandleConnection: Trouble with Accept: %q", err.Error())
		return
	}
	conn.SetReadLimit(h.readLimit)
	defer conn.Close(websocket.StatusNormalClosure, "Normal closure")

	h.wsService.StartService(ctx, conn, userId)
}
