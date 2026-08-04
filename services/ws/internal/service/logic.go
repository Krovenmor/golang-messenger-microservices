package service

import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type wsService struct {
	sub Subscriber
	pub Publisher

	connCounter atomic.Int64
	connMap     sync.Map
}

func NewWsService(sub Subscriber, pub Publisher) *wsService {
	return &wsService{
		sub: sub,
		pub: pub,
	}
}

func (s *wsService) startLog(userId uuid.UUID) func() {
	currentClients := s.connCounter.Add(1)
	val, isLoaded := s.connMap.Load(userId)
	iVal := 0
	if isLoaded {
		iVal = val.(int)
	}
	iVal++
	s.connMap.Store(userId, iVal)
	log.Printf("Сlient connected: Clients counter:%d, Client number of alive conns:%d, UserID: %s", currentClients, iVal, userId.String())

	return func() {
		currentClients := s.connCounter.Add(-1)
		val, _ := s.connMap.Load(userId)
		iVal := val.(int) - 1
		s.connMap.Store(userId, iVal)
		log.Printf("Сlient disconnected: Clients counter:%d, Client number of alive conns:%d, UserID: %s", currentClients, iVal, userId.String())
	}
}

func (s *wsService) StartService(ctx context.Context, conn *websocket.Conn, userId uuid.UUID) {
	chUserEvents, cancelEvents, err := s.sub.Subscribe(ctx, userId.String())
	if err != nil {
		log.Printf("StartService: Trouble with Subscribe, err: %q", err.Error())
		conn.Close(websocket.StatusInternalError, "failed to subscribe")
		return
	}
	defer cancelEvents()

	endLog := s.startLog(userId)
	defer endLog()

	cCtx, cancelCtx := context.WithCancel(ctx)
	worker := newWsWorker(s.sub, s.pub, conn, cCtx, cancelCtx, userId)
	go worker.startReader()
	worker.startWriter(chUserEvents)
}
