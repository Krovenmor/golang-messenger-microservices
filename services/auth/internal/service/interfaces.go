package service

import (
	j "MyMessenger/pkg/jwt"
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepo interface {
	AddNewUser(ctx context.Context, userId uuid.UUID, login, password string) error
	GetUser(ctx context.Context, login string) (uuid.UUID, string, error)
	IsUserExists(ctx context.Context, login string) error
	GetUserInfo(ctx context.Context, userId uuid.UUID) (*UserInfo, error)

	SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error
	FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error

	DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error
	DeleteExpiredRefreshTokens(ctx context.Context, userId uuid.UUID) error
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
