package service

import (
	"MyMessenger/services/profile/internal/config"
	"context"
	"log"

	"github.com/google/uuid"
)

type profileService struct {
	repo ProfileRepo
	pub  ProfilePublisher
	conf *config.ProfileConfig
}

func NewProfileService(repo ProfileRepo, pub ProfilePublisher, conf *config.ProfileConfig) *profileService {
	return &profileService{repo: repo, pub: pub, conf: conf}
}

func (s *profileService) NewProfile(ctx context.Context, profile Profile) error {
	name, err := s.checkName(profile.Name)
	if err != nil {
		return err
	}
	profile.Name = name

	err = s.repo.NewProfile(ctx, &profile)
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

func (s *profileService) GetProfilesById(ctx context.Context, userIds []uuid.UUID) ([]Profile, error) {
	return s.repo.GetProfilesById(ctx, userIds)
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

func (s *profileService) ChangeProfileName(ctx context.Context, userId uuid.UUID, name string) error {
	name, err := s.checkName(name)
	if err != nil {
		return err
	}
	return s.repo.UpdateName(ctx, userId, name)
}

func (s *profileService) ChangeProfileBio(ctx context.Context, userId uuid.UUID, bio string) error {
	bio, err := s.checkBio(bio)
	if err != nil {
		return err
	}
	return s.repo.UpdateBio(ctx, userId, bio)
}
