package service

import (
	"context"

	"github.com/google/uuid"
)

type ProfileRepo interface {
	NewProfile(ctx context.Context, profile *Profile) error
	DelProfile(ctx context.Context, profileId uuid.UUID) error

	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfilesById(ctx context.Context, userIds []uuid.UUID) ([]Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)

	AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error
	DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error

	UpdateName(ctx context.Context, userId uuid.UUID, name string) error
	UpdateBio(ctx context.Context, userId uuid.UUID, bio string) error
}

type ProfilePublisher interface {
	PubNewProfile(ctx context.Context, userId uuid.UUID) error
}
