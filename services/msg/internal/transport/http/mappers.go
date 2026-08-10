package http

import (
	"MyMessenger/services/msg/internal/service"

	"github.com/google/uuid"
)

func ToServiceProfile(p *ProfileBody, userId uuid.UUID) *service.Profile {
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

func FromServiceProfile(p *service.Profile) *ProfileBody {
	return &ProfileBody{
		UserId:     p.UserId,
		Name:       p.Name,
		UserName:   p.UserName,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
		CreatedAt:  p.CreatedAt,
		KeyNonce:   p.KeyNonce,
	}
}

func ToPublicProfileBody(p *service.Profile) *ProfilePublicBody {
	return &ProfilePublicBody{
		UserId:    p.UserId,
		Name:      p.Name,
		UserName:  p.UserName,
		PublicKey: p.PublicKey,
		CreatedAt: p.CreatedAt,
	}
}

func ToServiceMsg(m *PostMessageIncomeBody, userId uuid.UUID) *service.Message {
	return &service.Message{
		SenderId: userId,
		Message:  &m.Msg,
	}
}
