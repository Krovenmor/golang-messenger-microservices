package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthRepo interface {
	AddNewUser(ctx context.Context, userId uuid.UUID, login, password string) error
	GetUser(ctx context.Context, login, password string) (uuid.UUID, error)

	SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error
	FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error

	DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error
}
