package service

import (
	"time"

	"github.com/google/uuid"
)

type MediaTypes string
type MediaSubTypes string

const (
	PhotoType MediaTypes = "photo"

	AvatarSubType MediaSubTypes = "avatar"
)

type MediaInfo struct {
	MediaId  uuid.UUID
	Type     MediaTypes
	SubType  MediaSubTypes
	Size     int64
	IsPublic bool
	AddedAt  *time.Time
}

type ProfileInfo struct {
	MaxSpace    int
	SpaceFilled int
	FilesSaved  int
}

func ToAvatarMediaInfo(photoId uuid.UUID, size int64) *MediaInfo {
	return &MediaInfo{
		MediaId:  photoId,
		Type:     PhotoType,
		SubType:  AvatarSubType,
		Size:     size,
		IsPublic: true,
	}
}
