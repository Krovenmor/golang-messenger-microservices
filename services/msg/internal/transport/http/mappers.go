package http

import (
	"MyMessenger/services/msg/internal/service"

	"github.com/google/uuid"
)

func ToServiceProfile(p *ProfileBody, userId uuid.UUID) *service.Profile {
	return &service.Profile{
		UserId:     userId,
		Name:       p.Name,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
	}
}

func FromServiceProfile(p *service.Profile) *ProfileBody {
	return &ProfileBody{
		Name:       p.Name,
		PublicKey:  p.PublicKey,
		PrivateKey: p.PrivateKey,
		KDFSalt:    p.KDFSalt,
	}
}

func ToServiceMsg(m *PostMessageIncomeBody, userId uuid.UUID) *service.Message {
	return &service.Message{
		SenderId: userId,
		Message:  m.Msg,
	}
}
