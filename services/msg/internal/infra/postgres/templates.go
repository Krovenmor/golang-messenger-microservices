package postgres

import (
	"MyMessenger/services/msg/internal/service"
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
		&profile.CreatedAt,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return &profile, nil
}
