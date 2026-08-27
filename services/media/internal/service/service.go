package service

import (
	"context"
	"image"

	"github.com/google/uuid"
)

type MediaService interface {
	SaveAvatar(ctx context.Context, userId uuid.UUID, img image.Image) (uuid.UUID, error)
	DeleteAvatar(ctx context.Context, userId, photoId uuid.UUID) error

	GetProfileInfo(ctx context.Context, userId uuid.UUID) (*ProfileInfo, error)
	GetProfileMediaInfo(ctx context.Context, userId, fromId uuid.UUID, quantity int) ([]MediaInfo, error)
}
