package postgres

import (
	"MyMessenger/services/auth/internal/domain"
	"MyMessenger/services/auth/internal/infra/postgres/queries"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostagreRepo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) (*PostagreRepo, error) {
	if pool == nil {
		return nil, fmt.Errorf("NewRepo(): pool == nil")
	}
	return &PostagreRepo{pool: pool}, nil
}

func (p *PostagreRepo) AddNewUser(ctx context.Context, userId uuid.UUID, login, password, email string) error {
	c, err := p.pool.Exec(ctx, queries.AddUser, userId, login, password, email)
	if err != nil {
		return betterError(err)
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("can't add new user")
	}
	return nil
}

func (p *PostagreRepo) GetUser(ctx context.Context, login string) (uuid.UUID, string, error) {
	var userId uuid.UUID
	var password string
	err := p.pool.QueryRow(ctx, queries.GetUser, login).Scan(&userId, &password)
	if err != nil {
		return userId, password, betterError(err)
	}
	return userId, password, nil
}

func (p *PostagreRepo) GetUserTokens(ctx context.Context, userId uuid.UUID) (*domain.UserTokens, error) {
	var info domain.UserTokens
	err := p.pool.QueryRow(ctx, queries.GetUserTokens, userId).Scan(&info.Login, &info.RTokens)
	if err != nil {
		return nil, betterError(err)
	}
	return &info, nil
}

func (p *PostagreRepo) GetUserInfo(ctx context.Context, userId uuid.UUID) (*domain.UserInfo, error) {
	var info domain.UserInfo
	err := p.pool.QueryRow(ctx, queries.GetUserInfo, userId).Scan(&info.Login, &info.Email)
	if err != nil {
		return nil, betterError(err)
	}
	return &info, nil
}

func (p *PostagreRepo) IsUserExists(ctx context.Context, login string) error {
	var isExists int
	err := p.pool.QueryRow(ctx, queries.CheckUserExists, login).Scan(&isExists)
	if err != nil {
		return betterError(err)
	}
	if isExists != 1 {
		return domain.ErrNotFound
	}
	return nil
}

func (p *PostagreRepo) SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error {
	c, err := p.pool.Exec(ctx, queries.SaveRefresh, userId, rToken, expAt)
	if err != nil {
		return betterError(err)
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("can't save refresh")
	}
	return nil
}

func (p *PostagreRepo) DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	c, err := p.pool.Exec(ctx, queries.DeleteRefresh, userId, rToken)
	if err != nil {
		return betterError(err)
	}
	if c.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (p *PostagreRepo) DeleteExpiredRefreshTokens(ctx context.Context, userId uuid.UUID) error {
	_, err := p.pool.Exec(ctx, queries.ClearExpiredRefreshTokens, userId)
	if err != nil {
		return betterError(err)
	}
	return nil
}
