package ws

import (
	"MyMessenger/pkg/broker"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WSHandler struct {
	sub Subscriber
}

func NewWSHandler(subscriber Subscriber) *WSHandler {
	return &WSHandler{
		sub: subscriber,
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
	defer conn.Close(websocket.StatusNormalClosure, "Normal closure")

	log.Printf("Client has arrived, Address:%q, ID:%q", r.RemoteAddr, userId.String())

	ch, cancel, err := h.sub.Subscribe(ctx, fmt.Sprintf("user:%s:events", userId.String()))
	if err != nil {
		log.Printf("HandleConnection: Trouble with Subscribe, err: %q", err.Error())
		conn.Close(websocket.StatusInternalError, "failed to subscribe")
		return
	}
	defer cancel()

	go func() {
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				log.Printf("conn.Read(), UserID: %s, err: %v", userId.String(), err)
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client disconnected (context canceled), UserID: %s", userId.String())
			return

		case msg, ok := <-ch:
			if !ok {
				log.Printf("Event channel closed for UserID: %s", userId.String())
				return
			}

			var event broker.Event
			if err := json.Unmarshal(msg, &event); err != nil {
				log.Printf("Trouble with Unmarshaling broker data: %v", err)
				continue
			}

			err = conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				log.Printf("Client write failed (likely disconnected), UserID: %s, err: %v", userId.String(), err)
				return
			}

		case <-ticker.C:
			ctxWithTime, cancelCTX := context.WithTimeout(ctx, time.Second*3)
			err := conn.Ping(ctxWithTime)
			cancelCTX()
			if err != nil {
				log.Printf("Client ping failed (likely disconnected), UserID: %s, err: %v", userId.String(), err)
				return
			}
		}
	}
}
