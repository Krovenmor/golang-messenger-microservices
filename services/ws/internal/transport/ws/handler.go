package ws

import (
	"MyMessenger/pkg/broker"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type WSHandler struct {
	sub         Subscriber
	connCounter atomic.Int64
	connMap     sync.Map
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

	ch, cancel, err := h.sub.Subscribe(ctx, fmt.Sprintf("user:%s:events", userId.String()))
	if err != nil {
		log.Printf("HandleConnection: Trouble with Subscribe, err: %q", err.Error())
		conn.Close(websocket.StatusInternalError, "failed to subscribe")
		return
	}
	defer cancel()

	logError := "No error"
	startLog := func() func() {
		currentClients := h.connCounter.Add(1)
		val, isLoaded := h.connMap.Load(userId)
		var iVal int
		if !isLoaded {
			iVal = 0
		} else {
			iVal = val.(int)
		}
		iVal++
		h.connMap.Store(userId, iVal)
		log.Printf("Сlient connected: Clients counter:%d, Client number of alive conns:%d, UserID: %s", currentClients, iVal, userId.String())

		return func() {
			currentClients := h.connCounter.Add(-1)
			val, _ := h.connMap.Load(userId)
			iVal := val.(int) - 1
			log.Printf("Сlient disconnected: Clients counter:%d, Client number of alive conns:%d, UserID: %s, err: %q", currentClients, iVal, userId.String(), logError)
			h.connMap.Store(userId, iVal)
		}
	}
	endLog := startLog()
	defer endLog()

	go func() {
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				logError = fmt.Sprintf("conn.Read(), err: %v", err)
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logError = "context canceled"
			return

		case msg, ok := <-ch:
			if !ok {
				logError = "Event channel closed"
				return
			}

			var event broker.Event
			if err := json.Unmarshal(msg, &event); err != nil {
				log.Printf("Trouble with Unmarshaling broker data: %v", err)
				continue
			}

			err = conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				logError = fmt.Sprintf("Client write failed, err: %v", err)
				return
			}

		case <-ticker.C:
			ctxWithTime, cancelCTX := context.WithTimeout(ctx, time.Second*3)
			err := conn.Ping(ctxWithTime)
			cancelCTX()
			if err != nil {
				logError = fmt.Sprintf("Client ping failed, err: %v", err)
				return
			}
		}
	}
}
