package infra

import (
	"MyMessenger/services/auth/internal/infra/queries"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyExists     = errors.New("already exists")
	ErrAlreadyExistsUser = errors.New("user already exists")
	ErrNotFound          = errors.New("not found")
)

type PostagreRepo struct {
	pool *pgxpool.Pool
	q    queries.Queries
}

func betterError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "users_login_key":
				return ErrAlreadyExistsUser
			}
		}
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
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

func (p *PostagreRepo) GetUser(ctx context.Context, login, password string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := p.pool.QueryRow(ctx, p.q.GetUser, login, password).Scan(&userId)
	if err != nil {
		return userId, betterError(err)
	}
	return userId, nil
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
