package postgres

import (
	"MyMessenger/pkg/repo"
	"MyMessenger/services/profile/internal/infra/postgres/queries"
	"MyMessenger/services/profile/internal/service"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *PostagreRepo) NewProfile(ctx context.Context, profile *service.Profile) error {
	_, err := r.pool.Exec(ctx, r.q.NewProfile,
		profile.UserId,
		profile.UserName,
		profile.Name,
		profile.PublicKey,
		profile.PrivateKey,
		profile.KDFSalt,
		profile.KeyNonce,
	)
	if err != nil {
		return getErrorMsg(err)
	}
	return nil
}

func (r *PostagreRepo) DelProfile(ctx context.Context, profileId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.DelProfile, profileId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}

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

func (r *PostagreRepo) GetProfileById(ctx context.Context, userId uuid.UUID) (*service.Profile, error) {
	return getProfileByVal(ctx, r.pool, r.q.GetProfileId, userId)
}

func (r *PostagreRepo) GetProfilesById(ctx context.Context, userIds []uuid.UUID) ([]service.Profile, error) {
	sl, err := repo.GetSliceQueryByPos[service.Profile](ctx, r.pool, r.q.GetProfilesId, userIds)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []service.Profile{}, nil
		}
		return sl, getErrorMsg(err)
	}
	return sl, nil
}

func (r *PostagreRepo) GetProfileByUserName(ctx context.Context, username string) (*service.Profile, error) {
	return getProfileByVal(ctx, r.pool, r.q.GetProfileUserName, username)
}

func (r *PostagreRepo) AddAvatarToProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.AddAvatarPhoto, userId, avatarId.String())
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (r *PostagreRepo) DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.DelAvatarPhoto, userId, avatarId.String())
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}
