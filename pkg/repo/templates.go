package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetQueryByPos[R any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) (*R, error) {
	var rVal R
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rVal, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[R])
	if err != nil {
		return nil, err
	}

	return &rVal, nil
}

func GetSliceQueryByPos[R any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) ([]R, error) {
	var rVal []R
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return rVal, err
	}
	defer rows.Close()

	rVal, err = pgx.CollectRows(rows, pgx.RowToStructByPos[R])
	if err != nil {
		return nil, err
	}

	return rVal, nil
}

func GetSliceQueryByType[R any](ctx context.Context, pool *pgxpool.Pool, query string, args ...any) ([]R, error) {
	var rVal []R
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return rVal, err
	}
	defer rows.Close()

	rVal, err = pgx.CollectRows(rows, pgx.RowTo[R])
	if err != nil {
		return nil, err
	}

	return rVal, nil
}

func GetSliceQueryByFunc[R any](ctx context.Context, pool *pgxpool.Pool, query string,
	collectFunc func(row pgx.CollectableRow) (R, error),
	args ...any) ([]R, error) {

	var rVal []R
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return rVal, err
	}
	defer rows.Close()

	rVal, err = pgx.CollectRows(rows, pgx.RowToFunc[R](collectFunc))
	if err != nil {
		return nil, err
	}

	return rVal, nil
}
