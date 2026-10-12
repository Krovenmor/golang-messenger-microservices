package service

import (
	"MyMessenger/services/auth/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	StartRegistration(ctx context.Context, login, password, email string) (uuid.UUID, time.Duration, error)
	ResendCode(ctx context.Context, registrationId uuid.UUID) (time.Duration, error)
	ConfirmRegistration(ctx context.Context, registrationId uuid.UUID, code string) error

	LogIn(ctx context.Context, login, password string) (*domain.Tokens, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error)

	UpdateTokens(ctx context.Context, rToken string) (*domain.Tokens, error)

	BanUser(ctx context.Context, userId uuid.UUID, ttl time.Duration) error
}
