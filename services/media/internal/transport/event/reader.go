package event

import (
	"MyMessenger/pkg/broker"
	"MyMessenger/services/media/internal/service"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

const (
	serviceName = "media-service"
)

type ProfileReader struct {
	media service.MediaService
}

func NewProfileReader(media service.MediaService) *ProfileReader {
	return &ProfileReader{media: media}
}

func (r *ProfileReader) GetServiceName() string {
	return serviceName
}

func (r *ProfileReader) OnNewProfileEvent(ctx context.Context, event broker.ProfilePayload) error {
	id, err := uuid.Parse(event.UserId)
	if err != nil {
		return err
	}
	err = r.media.NewProfile(ctx, id)
	if err != nil {
		if !errors.Is(err, service.ErrAlreadyExists) {
			return err
		}
		log.Printf("ProfileReader.OnNewProfileEvent(), profile AlreadyExists, profile: %q", id)
	}
	return nil
}
