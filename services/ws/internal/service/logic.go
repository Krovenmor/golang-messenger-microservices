package service

import (
	"MyMessenger/services/ws/internal/config"
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type wsService struct {
	sub       Subscriber
	pub       Publisher
	msgClient MessageClient

	connCounter atomic.Int64
	connMap     sync.Map

	workerConf *config.WsConfig
}

func NewWsService(sub Subscriber, pub Publisher, msgClient MessageClient, conf *config.WsConfig) *wsService {
	return &wsService{
		sub:        sub,
		pub:        pub,
		msgClient:  msgClient,
		workerConf: conf,
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

func (s *wsService) StartService(ctx context.Context, conn Connector, userId uuid.UUID, accessToken string) error {
	endLog := s.startLog(userId)
	defer endLog()

	worker := newWsWorker(s.sub, s.pub, s.msgClient, conn, s.workerConf)
	return worker.startAll(ctx, accessToken, userId)
}
