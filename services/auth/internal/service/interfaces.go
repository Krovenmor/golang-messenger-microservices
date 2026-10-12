package service

import (
	j "MyMessenger/pkg/jwt"
	"MyMessenger/services/auth/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepo interface {
	AddNewUser(ctx context.Context, userId uuid.UUID, login, password, email string) error
	GetUser(ctx context.Context, login string) (uuid.UUID, string, error)
	CheckTaken(ctx context.Context, login, email string) (bool, bool, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error)

	GetUserTokens(ctx context.Context, userId uuid.UUID) (*domain.UserTokens, error)
	SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error

	DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error
	DeleteExpiredRefreshTokens(ctx context.Context, userId uuid.UUID) error
}

type RegistrationStore interface {
	Save(ctx context.Context, reg *domain.Registration, ttl time.Duration) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Registration, error)
	Delete(ctx context.Context, id uuid.UUID) error

	SaveCode(ctx context.Context, id uuid.UUID, code string, ttl time.Duration) error
	GetCode(ctx context.Context, id uuid.UUID) (string, error)

	StartCooldown(ctx context.Context, email string, ttl time.Duration) error
	StopCooldown(ctx context.Context, email string) error
}

type Publisher interface {
	PublishEmailVerification(ctx context.Context, email, code string) error
	PublishRegistrationAttempt(ctx context.Context, email string) error
}

type BanStore interface {
	BanLogin(ctx context.Context, login string, ttl time.Duration) error
	BanTokens(ctx context.Context, rTokenHashes []string, ttl time.Duration) error
}

type AuthHasher interface {
	Hash(password string) (string, error)
	Comp(password, hash string) error
}

type TokenGenerator interface {
	GenToken(userId uuid.UUID, tTl time.Duration, tType j.TokenType) (string, error)
}

type TokenChecker interface {
	IsValidToken(cToken string) (*j.TokenClaims, error)
}
