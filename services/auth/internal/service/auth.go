package service

import (
	"context"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	Login   string
	RTokens []string
}

type AuthService interface {
	Register(ctx context.Context, login, password string) error
	LogIn(ctx context.Context, login, password string) (*Tokens, error)
	UpdateTokens(ctx context.Context, rToken string) (*Tokens, error)
	GetInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error)
}
