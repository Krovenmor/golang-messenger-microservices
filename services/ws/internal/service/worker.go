package service

import (
	"MyMessenger/pkg/broker"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

var StatusOk json.RawMessage
var ErrBadJson json.RawMessage
var ErrBadStatus json.RawMessage

func ComputeJsonResponses() {
	StatusOk, _ = json.Marshal(Response{
		Code: http.StatusOK,
		Msg:  "OK",
	})
	ErrBadJson, _ = json.Marshal(Response{
		Code: http.StatusBadRequest,
		Msg:  "Bad JSON",
	})
	ErrBadStatus, _ = json.Marshal(Response{
		Code: http.StatusBadRequest,
		Msg:  "Bad Status",
	})
}

type wsWorker struct {
	cChannel chan []byte
	ctx      context.Context
	sub      Subscriber
	pub      Publisher
	conn     *websocket.Conn
	cancel   context.CancelFunc

	mutex             sync.Mutex
	userId            string
	currentUserStatus broker.Status
}

func newWsWorker(sub Subscriber, pub Publisher, conn *websocket.Conn, ctx context.Context, cancel context.CancelFunc, userId uuid.UUID) *wsWorker {
	worker := &wsWorker{
		sub:    sub,
		pub:    pub,
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,

		userId:            userId.String(),
		currentUserStatus: broker.Offline,

		cChannel: make(chan []byte, 2),
	}
	return worker
}

func (s *wsWorker) changeStatus(newStatus broker.Status) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.currentUserStatus == newStatus {
		return
	}

	ctx := s.ctx
	if ctx.Err() != nil {
		log.Printf("context gone...")
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
	}

	log.Printf("publishing new status: %d", newStatus)
	s.pub.PublishUserStatus(ctx, broker.StatusEvent{
		UserId:    s.userId,
		Status:    newStatus,
		EventTime: time.Now().Unix(),
	})

	s.currentUserStatus = newStatus
}

func (s *wsWorker) isValidStatus(status broker.Status) bool {
	return broker.IsValidStatus(status) && status != broker.Offline
}

func (s *wsWorker) startReader() {
	log.Printf("worker reader on")
	defer log.Printf("worker reader off")
	defer s.cancel()

	for {
		t, data, err := s.conn.Read(s.ctx)
		if t == websocket.MessageText {
			var event NewStatusEvent
			unmarshalErr := json.Unmarshal(data, &event)
			if unmarshalErr == nil {
				status := broker.Status(event.NewStatus)
				if s.isValidStatus(status) {
					s.changeStatus(status)
					s.cChannel <- StatusOk
				} else {
					log.Printf("Trouble with data from user, not valid status %d", status)
					s.cChannel <- ErrBadStatus
				}
			} else {
				log.Printf("Trouble with unmarshaling data from user, err: %q", unmarshalErr)
				s.cChannel <- ErrBadJson
			}
		}
		if err != nil {
			return
		}
	}
}

func (s *wsWorker) startWriter(chEvents <-chan []byte) {
	log.Printf("worker writer on")
	defer log.Printf("worker writer off")

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	s.changeStatus(broker.Online)
	defer s.changeStatus(broker.Offline)
	defer s.cancel()

	for {
		select {
		case <-s.ctx.Done():
			return

		case msg, ok := <-chEvents:
			if !ok {
				return
			}
			err := s.conn.Write(s.ctx, websocket.MessageText, msg)
			if err != nil {
				return
			}

		case errMsg, ok := <-s.cChannel:
			if !ok {
				log.Printf("Trouble with error channel")
				return
			}
			err := s.conn.Write(s.ctx, websocket.MessageText, errMsg)
			if err != nil {
				return
			}

		case <-ticker.C:
			ctxWithTime, cancelCTX := context.WithTimeout(s.ctx, time.Second*3)
			err := s.conn.Ping(ctxWithTime)
			cancelCTX()
			if err != nil {
				return
			}
		}
	}
}
