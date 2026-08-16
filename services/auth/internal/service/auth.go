package service

import (
	"context"

	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type UserTokens struct {
	Login   string
	RTokens []string
}

type UserInfo struct {
	Login string
	Email string
}

type AuthService interface {
	SendCodeEmail(ctx context.Context, email string) (int, int, error)

	Register(ctx context.Context, login, password, email, code string) error
	LogIn(ctx context.Context, login, password string) (*Tokens, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error)

	UpdateTokens(ctx context.Context, rToken string) (*Tokens, error)
	GetUserTokens(ctx context.Context, userId uuid.UUID) (*UserTokens, error)
}
