package service

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type profileService struct {
	repo ProfileRepo
	pub  ProfilePublisher
}

func NewProfileService(repo ProfileRepo, pub ProfilePublisher) *profileService {
	return &profileService{repo: repo, pub: pub}
}

func (s *profileService) NewProfile(ctx context.Context, profile Profile) error {
	err := s.repo.NewProfile(ctx, &profile)
	if err != nil {
		return err
	}
	err = s.pub.PubNewProfile(ctx, profile.UserId)
	if err != nil {
		log.Printf("profileService: trouble with PubNewProfile, err: %q", err)
		delErr := s.repo.DelProfile(ctx, profile.UserId)
		if delErr != nil {
			log.Printf("profileService: trouble with DelProfile, err: %q", delErr)
		}
		return ErrInternal
	}
	return nil
}

func (s *profileService) GetProfileById(ctx context.Context, userId uuid.UUID) (*Profile, error) {
	return s.repo.GetProfileById(ctx, userId)
}

func (s *profileService) GetProfileByUserName(ctx context.Context, username string) (*Profile, error) {
	return s.repo.GetProfileByUserName(ctx, username)
}

func (s *profileService) AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error {
	return s.repo.AddAvatarToProfile(ctx, userId, avatarId)
}

func (s *profileService) DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error {
	return s.repo.DelAvatarFromProfile(ctx, userId, avatarId)
}
