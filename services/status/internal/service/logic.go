package service

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/services/status/internal/config"
	"context"
	"errors"
	"log"
)

var (
	ErrWrongStatus   = errors.New("wrong status")
	ErrSavingTrouble = errors.New("trouble with saving")
)

var defaultUserStatus = UserStatus{Status: broker.Offline, LastSeen: 0}

type statusService struct {
	repo StatusRepo
	conf *config.ServiceConfig
}

func NewStatusService(repo StatusRepo, conf *config.ServiceConfig) *statusService {
	return &statusService{repo: repo, conf: conf}
}

func (s *statusService) GetStatus(ctx context.Context, userId string) UserStatus {
	status, err := s.repo.GetStatus(ctx, userId)
	if err != nil {
		log.Printf("Error with repo.GetStatus(ctx, userId), err: %q", err.Error())
		return defaultUserStatus
	}
	return *status
}

func (s *statusService) SaveStatus(ctx context.Context, event broker.StatusPayload) error {
	if !broker.IsValidStatus(event.Status) {
		return ErrWrongStatus
	}
	err := s.repo.SaveStatus(
		ctx,
		event,
		s.conf.EntriesTTL,
	)
	if err != nil {
		log.Printf("Trouble with saving %q", err.Error())
		return ErrSavingTrouble
	}
	return nil
}
