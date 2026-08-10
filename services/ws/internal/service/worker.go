package service

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/services/ws/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

const (
	bufferLength     = 10
	respBufferLength = 2
)

type wsWorker struct {
	eventsMutex   sync.Mutex
	eventsChannel chan []byte
	eventsClose   func()

	respWriteChannel chan []byte

	ctx       context.Context
	ctxCancel context.CancelFunc
	sub       Subscriber
	pub       Publisher
	conn      Connector
	msgClient MessageClient

	userCancel context.CancelFunc

	userId            uuid.UUID
	currentUserStatus broker.Status

	conf *config.WsConfig
}

func newWsWorker(sub Subscriber, pub Publisher, msgClient MessageClient, conn Connector, conf *config.WsConfig) *wsWorker {

	eventsChannel := make(chan []byte, bufferLength)

	return &wsWorker{
		sub:       sub,
		pub:       pub,
		conn:      conn,
		msgClient: msgClient,

		currentUserStatus: broker.Offline,

		respWriteChannel: make(chan []byte, respBufferLength),
		eventsChannel:    eventsChannel,
		eventsClose:      func() { close(eventsChannel) },

		conf: conf,
	}
}

func (s *wsWorker) writeToMain(ch <-chan []byte, cancel func()) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}

			s.eventsMutex.Lock()
			select {
			case s.eventsChannel <- msg:
			default:
				log.Printf("writeToMain: eventsChannel full!")
			}
			s.eventsMutex.Unlock()

		case <-s.ctx.Done():
			cancel()
			return
		}
	}
}

func (s *wsWorker) subAll(accessToken string) error {
	userString := s.userId.String()

	ctx, caneclCtx := context.WithTimeout(s.ctx, s.conf.ContextWaitTime)
	chats, err := s.msgClient.GetAllUserChats(ctx, userString, accessToken)
	caneclCtx()

	if err != nil {
		log.Printf("trouble with GetAllUserChats, err: %q", err)
		return err
	}

	chUserEvents, cancelEvents, err := s.sub.SubscribeOnUserEvents(s.ctx, userString)
	if err != nil {
		log.Printf("trouble with Subscribe, err: %q", err)
		return ErrInternal
	}
	go s.writeToMain(chUserEvents, cancelEvents)

	for _, chat := range chats {
		chatCh, chatCancel, err := s.sub.SubscribeOnChatEvents(s.ctx, chat.String())
		if err != nil {
			log.Printf("trouble with SubscribeOnChatEvents, err: %q, chatId: %q", err, chat)
			return ErrInternal
		}
		go s.writeToMain(chatCh, chatCancel)
	}
	return nil
}

func (s *wsWorker) changeStatus(newStatus broker.Status) {
	if s.currentUserStatus == newStatus {
		return
	}

	ctx := s.ctx
	if ctx.Err() != nil {
		log.Printf("changeStatus: context gone...")
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.conf.ContextWaitTime)
		defer cancel()
	}

	s.pub.PublishUserStatus(ctx, broker.StatusPayload{
		UserId:    s.userId.String(),
		Status:    newStatus,
		EventTime: time.Now().Unix(),
	})

	s.currentUserStatus = newStatus
}

func (s *wsWorker) isValidStatus(status broker.Status) bool {
	return broker.IsValidStatus(status) && status != broker.Offline
}

func unmarshalJson[T any](msg json.RawMessage) (T, error) {
	var body T
	dec := json.NewDecoder(bytes.NewReader(msg))
	dec.DisallowUnknownFields()
	err := dec.Decode(&body)
	if err != nil {
		log.Printf("Trouble with unmarshaling data from user, err: %q", err)
		return body, err
	}
	return body, nil
}

func (s *wsWorker) sendResponse(data []byte) bool {
	select {
	case s.respWriteChannel <- data:
		return true
	case <-s.ctx.Done():
		return false
	}
}

func (s *wsWorker) handleNewStatus(payload RPNewStatus) {
	status := broker.Status(payload.NewStatus)
	if s.isValidStatus(status) {
		s.changeStatus(status)
		s.sendResponse(StatusOkResp)
	} else {
		log.Printf("Trouble with data from user, not valid status %d", status)
		s.sendResponse(ErrBadStatusResp)
	}
}

func (s *wsWorker) handleNewUserSub(payload RPSubUser) {
	ctx, cancel := context.WithCancel(s.ctx)
	uCh, uCancel, err := s.sub.SubscribeOnUserStatuses(ctx, payload.UserId.String())

	if err != nil {
		log.Printf("handleNewUserSub: Trouble with SubscribeOnUserEvents, err: %q", err)
		s.sendResponse(ErrInternalResp)
		cancel()
		return
	}

	if s.userCancel != nil {
		s.userCancel()
	}

	s.userCancel = cancel

	go s.writeToMain(uCh, uCancel)
	s.sendResponse(StatusOkResp)
}

func (s *wsWorker) handleRequest(req Request) {
	switch req.Req {
	case NewStatusReq:
		body, err := unmarshalJson[RPNewStatus](req.Payload)
		if err != nil {
			s.sendResponse(ErrBadJsonResp)
			return
		}
		s.handleNewStatus(body)
	case SubUserReq:
		body, err := unmarshalJson[RPSubUser](req.Payload)
		if err != nil {
			s.sendResponse(ErrBadJsonResp)
			return
		}
		s.handleNewUserSub(body)
	default:
		s.sendResponse(ErrBadRequestResp)
		return
	}
}

func (s *wsWorker) isNormalError(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrConnNormalClosure) {
		return true
	}
	return false
}

func (s *wsWorker) banUser() {
	s.pub.PublishBanEvent(s.ctx, broker.BanRequestPayload{
		UserId: s.userId,
		Reason: broker.TooManyRequests,
	})
}

func (s *wsWorker) startReader() error {
	defer s.ctxCancel()

	limiter := rate.NewLimiter(rate.Limit(s.conf.LimitRate), s.conf.LimitBurst)
	violations := 0

	for {
		t, data, err := s.conn.Read(s.ctx)
		if err != nil {
			return err
		}
		if !limiter.Allow() {
			violations++
			if violations > s.conf.LimitViolations {
				s.banUser()
				return ErrTooManyRequests
			}
		}
		if t == TextType {
			req, unmarshalErr := unmarshalJson[Request](data)
			if unmarshalErr == nil {

				s.handleRequest(req)

			} else {
				log.Printf("Trouble with unmarshaling data from user, err: %q", unmarshalErr)

				if !s.sendResponse(ErrBadJsonResp) {
					return s.ctx.Err()
				}
			}
		}
	}
}

func (s *wsWorker) writeMessage(msg []byte) error {
	err := s.conn.Write(s.ctx, TextType, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *wsWorker) startWriter() error {
	defer s.ctxCancel()

	ticker := time.NewTicker(s.conf.TickerTiming)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()

		case msg, ok := <-s.eventsChannel:
			if !ok {
				log.Printf("Trouble with chatEventsChannel")
				return ErrInternal
			}
			err := s.writeMessage(msg)
			if err != nil {
				return err
			}

		case resp, ok := <-s.respWriteChannel:
			if !ok {
				log.Printf("Trouble with respWriteChannel")
				return ErrInternal
			}
			err := s.writeMessage(resp)
			if err != nil {
				return err
			}

		case <-ticker.C:
			ctxWithTime, cancelCTX := context.WithTimeout(s.ctx, s.conf.ContextWaitTime)
			err := s.conn.Ping(ctxWithTime)
			cancelCTX()
			if err != nil {
				if errors.Is(s.ctx.Err(), context.Canceled) {
					return s.ctx.Err()
				}
				log.Printf("ping failed for user %q: %v", s.userId, err)
				return err
			}
		}
	}
}

func (s *wsWorker) startAll(ctx context.Context, accessToken string, userId uuid.UUID) error {
	s.userId = userId
	ctxL, cancel := context.WithCancel(ctx)
	defer cancel()
	defer s.eventsClose()

	s.ctx = ctxL
	s.ctxCancel = cancel

	err := s.subAll(accessToken)
	if err != nil {
		log.Printf("Trouble with startAll, err: %q", err)
		return err
	}
	defer func() {
		if s.userCancel != nil {
			s.userCancel()
		}
	}()

	s.changeStatus(broker.Online)
	defer s.changeStatus(broker.Offline)

	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		collectedErrs error
	)
	addErr := func(err error) {
		if err == nil || s.isNormalError(err) {
			return
		}

		mu.Lock()
		collectedErrs = errors.Join(collectedErrs, err)
		mu.Unlock()

		s.ctxCancel()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.startReader()
		if err != nil {
			addErr(fmt.Errorf("error in reader: %w", err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.startWriter()
		if err != nil {
			addErr(fmt.Errorf("error in writer: %w", err))
		}
	}()

	wg.Wait()

	if collectedErrs != nil {
		log.Printf("Worker finished with multiple errors for user %s:\n%v", s.userId, collectedErrs)
		return collectedErrs
	}

	return nil
}
