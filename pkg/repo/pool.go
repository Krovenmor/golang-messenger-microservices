package repo

import (
	"MyMessenger/pkg/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

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
