package service

import (
	j "MyMessenger/pkg/jwt"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepo interface {
	AddNewUser(ctx context.Context, userId uuid.UUID, login, password, email string) error
	GetUser(ctx context.Context, login string) (uuid.UUID, string, error)
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error)
	IsUserExists(ctx context.Context, login string) error

	GetUserTokens(ctx context.Context, userId uuid.UUID) (*UserTokens, error)
	SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error
	FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error

	DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error
	DeleteExpiredRefreshTokens(ctx context.Context, userId uuid.UUID) error
}

type AuthTTLCache interface {
	Put(ctx context.Context, email, code string, ttl time.Duration) error
	Get(ctx context.Context, email string) (string, error)
	Del(ctx context.Context, email string) error
}

type Publisher interface {
	PublishEmailVerification(ctx context.Context, email, code string) error
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
