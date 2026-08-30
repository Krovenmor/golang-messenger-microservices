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
	return defExec(ctx, r.pool, r.q.DelProfile, profileId)
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
	return defExec(ctx, r.pool, r.q.AddAvatarPhoto, userId, avatarId.String())
}

func (r *PostagreRepo) DelAvatarFromProfile(ctx context.Context, userId uuid.UUID, avatarId uuid.UUID) error {
	return defExec(ctx, r.pool, r.q.DelAvatarPhoto, userId, avatarId.String())
}

func (r *PostagreRepo) UpdateBio(ctx context.Context, userId uuid.UUID, bio string) error {
	return defExec(ctx, r.pool, r.q.UpdateBio, userId, bio)
}

func (r *PostagreRepo) UpdateName(ctx context.Context, userId uuid.UUID, name string) error {
	return defExec(ctx, r.pool, r.q.UpdateName, userId, name)
}
