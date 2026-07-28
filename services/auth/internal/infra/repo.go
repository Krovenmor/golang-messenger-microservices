package infra

import (
	"MyMessenger/services/auth/internal/infra/queries"
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
		return err
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("Can't add new user")
	}
	return nil
}

func (p *PostagreRepo) GetUser(ctx context.Context, login, password string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := p.pool.QueryRow(ctx, p.q.GetUser, login, password).Scan(&userId)
	if err != nil {
		return userId, err
	}
	return userId, nil
}

func (p *PostagreRepo) SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error {
	c, err := p.pool.Exec(ctx, p.q.SaveRefresh, userId, rToken, expAt)
	if err != nil {
		return err
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("Can't save refresh")
	}
	return nil
}

func (p *PostagreRepo) FindRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	var isFound int
	err := p.pool.QueryRow(ctx, p.q.FindRefresh, userId, rToken).Scan(&isFound)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostagreRepo) DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	c, err := p.pool.Exec(ctx, p.q.DelRefresh, userId, rToken)
	if err != nil {
		return err
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("Can't del refresh")
	}
	return nil
}
