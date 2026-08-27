package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/media/internal/service"
	"time"

	"github.com/google/uuid"
)

type NewPhotoResponseBody struct {
	PhotoID string `json:"photoId"`
}

type MediaInfoResponseBody struct {
	MediaId  uuid.UUID  `json:"mediaId"`
	Type     string     `json:"type"`
	SubType  string     `json:"subType"`
	Size     int64      `json:"size"`
	IsPublic bool       `json:"isPublic"`
	AddedAt  *time.Time `json:"addedAt"`
}

type ProfileResponseBody struct {
	MaxSpace    int `json:"maxSpace"`
	SpaceFilled int `json:"spaceFilled"`
	SavedFiles  int `json:"savedFiles"`
}

func ToMediaInfoResponseBody(info service.MediaInfo) MediaInfoResponseBody {
	return MediaInfoResponseBody{
		MediaId:  info.MediaId,
		Type:     string(info.Type),
		SubType:  string(info.SubType),
		Size:     info.Size,
		IsPublic: info.IsPublic,
		AddedAt:  info.AddedAt,
	}
}

func ToMediaInfoSliceResponseBody(info []service.MediaInfo) []MediaInfoResponseBody {
	return utils.MapSlice(info, ToMediaInfoResponseBody)
}

func ToProfileResponseBody(info *service.ProfileInfo) *ProfileResponseBody {
	return &ProfileResponseBody{
		MaxSpace:    info.MaxSpace,
		SpaceFilled: info.SpaceFilled,
		SavedFiles:  info.FilesSaved,
	}
}
