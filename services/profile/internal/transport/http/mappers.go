package http

import (
	"MyMessenger/pkg/utils"
	"MyMessenger/services/profile/internal/service"

	"github.com/google/uuid"
)

func ToServiceProfile(p *NewProfileRequestBody, userId uuid.UUID) *service.Profile {
	return &service.Profile{
		UserId:     userId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
		KeyNonce:   p.KeyNonce,
	}
}

func ToPrivateProfileBody(p *service.Profile) *PrivateProfileResponseBody {
	return &PrivateProfileResponseBody{
		UserId:     p.UserId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
		CreatedAt:  p.CreatedAt,
		KeyNonce:   p.KeyNonce,
		Additional: p.Additional,
	}
}

func ToPublicProfileBody(p service.Profile) PublicProfileResponseBody {
	return PublicProfileResponseBody{
		UserId:     p.UserId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		CreatedAt:  p.CreatedAt,
		Additional: p.Additional,
	}
}

func ToPublicProfileBodies(p []service.Profile) []PublicProfileResponseBody {
	return utils.MapSlice(p, ToPublicProfileBody)
}
