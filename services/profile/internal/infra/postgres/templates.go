package postgres

import (
	"MyMessenger/services/profile/internal/service"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getProfileByVal[T any](ctx context.Context, pool *pgxpool.Pool, query string, val T) (*service.Profile, error) {
	var profile service.Profile
	err := pool.QueryRow(ctx, query, val).Scan(
		&profile.UserId,
		&profile.Name,
		&profile.UserName,
		&profile.PublicKey,
		&profile.PrivateKey,
		&profile.KDFSalt,
		&profile.KeyNonce,
		&profile.CreatedAt,
		&profile.Additional,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return &profile, nil
}

func defExec(ctx context.Context, pool *pgxpool.Pool, query string, args ...any) error {
	t, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}
