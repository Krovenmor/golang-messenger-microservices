package service

import (
	"context"

	"github.com/google/uuid"
)

type ProfileService interface {
	NewProfile(ctx context.Context, profile Profile) error

	GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error)
	GetProfilesById(ctx context.Context, userIds []uuid.UUID) ([]Profile, error)
	GetProfileByUserName(ctx context.Context, username string) (*Profile, error)

	AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error
	DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error

	ChangeProfileName(ctx context.Context, userId uuid.UUID, name string) error
	ChangeProfileBio(ctx context.Context, userId uuid.UUID, bio string) error
}
