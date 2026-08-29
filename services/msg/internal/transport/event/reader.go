package event

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/services/msg/internal/service"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

const (
	serviceName = "msg-service"
)

type ProfileReader struct {
	msg service.MessageService
}

func NewProfileReader(msg service.MessageService) *ProfileReader {
	return &ProfileReader{msg: msg}
}

func (r *ProfileReader) GetServiceName() string {
	return serviceName
}

func (r *ProfileReader) OnNewProfileEvent(ctx context.Context, event broker.ProfilePayload) error {
	id, err := uuid.Parse(event.UserId)
	if err != nil {
		return err
	}
	err = r.msg.NewProfile(ctx, id)
	if err != nil {
		if !errors.Is(err, service.ErrAlreadyExists) {
			return err
		}
		log.Printf("ProfileReader.OnNewProfileEvent(), profile AlreadyExists, profile: %q", id)
	}
	return nil
}
