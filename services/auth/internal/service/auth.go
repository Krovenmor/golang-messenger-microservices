package service

import (
	"context"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Register(ctx context.Context, login, password string) error
	LogIn(ctx context.Context, login, password string) (*Tokens, error)
	IsValidAccess(ctx context.Context, aToken string) error
	UpdateTokens(ctx context.Context, rToken string) (*Tokens, error)
}
