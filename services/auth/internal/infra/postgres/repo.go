package postgres

import (
	"MyMessenger/services/auth/internal/infra/postgres/queries"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostagreRepo struct {
	pool *pgxpool.Pool
	q    queries.Queries
}

func NewRepo(pool *pgxpool.Pool) (*PostagreRepo, error) {
	if pool == nil {
		return nil, fmt.Errorf("NewRepo(): pool == nil")
	}
	q, err := queries.GetQueries()
	if err != nil {
		return nil, err
	}
	return &PostagreRepo{
		pool: pool,
		q:    *q,
	}, nil
}

func (p *PostagreRepo) AddNewUser(ctx context.Context, userId uuid.UUID, login, password string) error {
	c, err := p.pool.Exec(ctx, p.q.AddUser, userId, login, password)
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
	err := p.pool.QueryRow(ctx, p.q.GetUser, login).Scan(&userId, &password)
	if err != nil {
		return userId, password, betterError(err)
	}
	return userId, password, nil
}

func (p *PostagreRepo) IsUserExists(ctx context.Context, login string) error {
	var isExists int
	err := p.pool.QueryRow(ctx, p.q.CheckUserExists, login).Scan(&isExists)
	if err != nil {
		return betterError(err)
	}
	if isExists != 1 {
		return ErrNotFound
	}
	return nil
}

func (p *PostagreRepo) SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error {
	c, err := p.pool.Exec(ctx, p.q.SaveRefresh, userId, rToken, expAt)
	if err != nil {
		return betterError(err)
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("can't save refresh")
	}
	return nil
}

func (p *PostagreRepo) FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	var isFound int
	err := p.pool.QueryRow(ctx, p.q.FindRefresh, userId, rToken).Scan(&isFound)
	if err != nil {
		return betterError(err)
	}
	return nil
}

func (p *PostagreRepo) DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	c, err := p.pool.Exec(ctx, p.q.DelRefresh, userId, rToken)
	if err != nil {
		return err
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("can't del refresh")
	}
	return nil
}

func (p *PostagreRepo) DeleteExpiredRefreshTokens(ctx context.Context, userId uuid.UUID) error {
	_, err := p.pool.Exec(ctx, p.q.ClrExpRefresh, userId)
	if err != nil {
		return err
	}
	return nil
}
