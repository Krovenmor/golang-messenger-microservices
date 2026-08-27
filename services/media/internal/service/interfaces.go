package service

import (
	"context"

	"github.com/google/uuid"
)

type MediaRepo interface {
	AddNewMedia(ctx context.Context, userId uuid.UUID, info *MediaInfo) error
	GetAvailableSpace(ctx context.Context, userId uuid.UUID) (int64, error)
	DelMedia(ctx context.Context, userId, mediaId uuid.UUID) error

	GetProfileInfo(ctx context.Context, userId uuid.UUID) (*ProfileInfo, error)
	GetProfileMediaInfo(ctx context.Context, userId, fromId uuid.UUID, quantity int) ([]MediaInfo, error)
}
