package infra

import (
	"MyMessenger/services/auth/internal/config"
	"MyMessenger/services/auth/internal/infra/queries"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type repoQueries struct {
	addUser     string
	getUser     string
	saveRefresh string
	findRefresh string
	delRefresh  string
}

type PostagreRepo struct {
	pool *pgxpool.Pool
	q    repoQueries
}

func NewPool(lf fx.Lifecycle, conf config.RepoConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), conf.ConnString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	lf.Append(fx.Hook{
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func getQueries() (*repoQueries, error) {
	addUser, err := queries.GetQuery("add_user.sql")
	if err != nil {
		return nil, err
	}
	getUser, err := queries.GetQuery("get_user.sql")
	if err != nil {
		return nil, err
	}
	saveRefresh, err := queries.GetQuery("save_refresh.sql")
	if err != nil {
		return nil, err
	}
	findRefresh, err := queries.GetQuery("find_refresh.sql")
	if err != nil {
		return nil, err
	}
	delRefresh, err := queries.GetQuery("delete_refresh.sql")
	if err != nil {
		return nil, err
	}
	return &repoQueries{
		addUser:     addUser,
		getUser:     getUser,
		saveRefresh: saveRefresh,
		findRefresh: findRefresh,
		delRefresh:  delRefresh,
	}, nil
}

func NewRepo(pool *pgxpool.Pool) (*PostagreRepo, error) {
	if pool == nil {
		return nil, fmt.Errorf("NewRepo(): pool == nil")
	}
	q, err := getQueries()
	if err != nil {
		return nil, err
	}
	return &PostagreRepo{
		pool: pool,
		q:    *q,
	}, nil
}

func (p *PostagreRepo) AddNewUser(ctx context.Context, userId uuid.UUID, login, password string) error {
	c, err := p.pool.Exec(ctx, p.q.addUser, userId, login, password)
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
	err := p.pool.QueryRow(ctx, p.q.getUser, login, password).Scan(&userId)
	if err != nil {
		return userId, err
	}
	return userId, nil
}

func (p *PostagreRepo) SaveRefresh(ctx context.Context, userId uuid.UUID, rToken string, expAt time.Time) error {
	c, err := p.pool.Exec(ctx, p.q.saveRefresh, userId, rToken, expAt)
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
	err := p.pool.QueryRow(ctx, p.q.findRefresh, userId, rToken).Scan(&isFound)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostagreRepo) DeleteRefresh(ctx context.Context, userId uuid.UUID, rToken string) error {
	c, err := p.pool.Exec(ctx, p.q.delRefresh, userId, rToken)
	if err != nil {
		return err
	}
	if c.RowsAffected() == 0 {
		return fmt.Errorf("Can't del refresh")
	}
	return nil
}
