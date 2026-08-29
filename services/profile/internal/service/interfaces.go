package service

import (
	"context"

	"github.com/google/uuid"
)

type ProfileRepo interface {
	NewProfile(ctx context.Context, profile *Profile) error
	DelProfile(ctx context.Context, profileId uuid.UUID) error

	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)

	AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error
	DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error
}

type ProfilePublisher interface {
	PubNewProfile(ctx context.Context, userId uuid.UUID) error
}
