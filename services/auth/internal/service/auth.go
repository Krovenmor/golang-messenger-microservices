package service

import (
	"MyMessenger/services/auth/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	SendCodeEmail(ctx context.Context, email string) (time.Duration, error)

	Register(ctx context.Context, login, password, email, code string) error
	LogIn(ctx context.Context, login, password string) (*domain.Tokens, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error)

	UpdateTokens(ctx context.Context, rToken string) (*domain.Tokens, error)

	BanUser(ctx context.Context, userId uuid.UUID, ttl time.Duration) error
}
