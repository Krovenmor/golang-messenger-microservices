package service

import (
	"MyMessenger/services/media/internal/config"
	"errors"
	"image/jpeg"

	"context"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/google/uuid"
)

const (
	avatarsDir = "avatars"

	defSpace    = 100 * 1024 * 1024
	defQuantity = 10

	minQuantity = 1
	maxQuantity = 100
)

type mediaSaver struct {
	avatarsPath string
	repo        MediaRepo
}

func NewMediaSaver(conf *config.MediaSaverConfig, repo MediaRepo) *mediaSaver {
	return &mediaSaver{
		repo:        repo,
		avatarsPath: conf.PublicPhotosSavingPath + avatarsDir + "/",
	}
}

func (s *mediaSaver) toFullPhotoPath(dir string, photoId uuid.UUID) string {
	return fmt.Sprintf("%s%s.jpeg", dir, photoId.String())
}

func (s *mediaSaver) saveImg(img image.Image, dir string, photoId uuid.UUID) (string, int64, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", -1, err
	}

	fullPath := s.toFullPhotoPath(dir, photoId)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", -1, err
	}
	defer file.Close()

	opts := &jpeg.Options{Quality: 70}
	err = jpeg.Encode(file, img, opts)
	if err != nil {
		os.Remove(fullPath)
		return "", -1, err
	}

	stat, _ := file.Stat()
	return fullPath, stat.Size(), nil
}

func (s *mediaSaver) delImg(filePath string) error {
	return os.Remove(filePath)
}

func (s *mediaSaver) getWorstSizeImg(img image.Image) int64 {
	return int64((img.Bounds().Dx() * img.Bounds().Dy() * 3) / 8)
}

func (s *mediaSaver) getSpace(ctx context.Context, userId uuid.UUID) (int64, error) {
	space, err := s.repo.GetAvailableSpace(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return defSpace, nil
		} else {
			return -1, err
		}
	}
	return space, nil
}

func (s *mediaSaver) NewProfile(ctx context.Context, userId uuid.UUID) error {
	return s.repo.NewProfile(ctx, userId)
}

func (s *mediaSaver) SaveAvatar(ctx context.Context, userId uuid.UUID, img image.Image) (uuid.UUID, error) {
	photoId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	space, err := s.getSpace(ctx, userId)
	if err != nil {
		return uuid.Nil, err
	}
	wSize := s.getWorstSizeImg(img)
	if space < wSize {
		return uuid.Nil, ErrNotEnoughSpace
	}

	fullPath, realSize, err := s.saveImg(img, s.avatarsPath, photoId)
	if err != nil {
		return uuid.Nil, err
	}

	if realSize > wSize {
		log.Printf("SaveAvatar: Bad formula for wSize!!!, realSize=%d, wSize=%d, img.Width=%d, img.Height=%d", realSize, wSize, img.Bounds().Dx(), img.Bounds().Dy())
	}

	avatarInfo := ToAvatarMediaInfo(photoId, realSize)
	err = s.repo.AddNewMedia(ctx, userId, avatarInfo)
	if err != nil {
		delErr := s.delImg(fullPath)
		if delErr != nil {
			log.Printf("SaveAvatar: IMAGE NOT DELETED, err: %q, filePath: %q", err, fullPath)
		}
		return uuid.Nil, err
	}

	return photoId, nil
}

func (s *mediaSaver) DeleteAvatar(ctx context.Context, userId, photoId uuid.UUID) error {
	err := s.repo.DelMedia(ctx, userId, photoId)
	if err != nil {
		return err
	}

	photoPath := s.toFullPhotoPath(s.avatarsPath, photoId)
	err = s.delImg(photoPath)
	if err != nil {
		log.Printf("DeleteAvatar: IMAGE NOT DELETED, err: %q, filePath: %q", err, photoPath)
		return ErrInteranl
	}
	return nil
}

func (s *mediaSaver) GetProfileInfo(ctx context.Context, userId uuid.UUID) (*ProfileInfo, error) {
	return s.repo.GetProfileInfo(ctx, userId)
}

func (s *mediaSaver) GetProfileMediaInfo(ctx context.Context, userId, fromId uuid.UUID, quantity int) ([]MediaInfo, error) {
	if fromId == uuid.Nil {
		fromId = uuid.Max
	}
	if quantity < minQuantity || quantity > maxQuantity {
		quantity = defQuantity
	}
	return s.repo.GetProfileMediaInfo(ctx, userId, fromId, quantity)
}
