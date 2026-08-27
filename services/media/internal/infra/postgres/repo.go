package postgres

import (
	"MyMessenger/pkg/repo"
	"MyMessenger/services/media/internal/infra/postgres/queries"
	"MyMessenger/services/media/internal/service"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
	q    *queries.Queries
}

func NewPostgresRepo(pool *pgxpool.Pool) (*PostgresRepo, error) {
	q, err := queries.GetQueries()
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{pool: pool, q: q}, nil
}

func (p *PostgresRepo) AddNewMedia(ctx context.Context, userId uuid.UUID, info *service.MediaInfo) error {
	t, err := p.pool.Exec(ctx, p.q.NewData, userId,
		info.MediaId, info.Type, info.SubType, info.Size, info.IsPublic,
	)
	if err != nil {
		return ToServiceError(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (p *PostgresRepo) GetAvailableSpace(ctx context.Context, userId uuid.UUID) (int64, error) {
	var aSpace int64
	err := p.pool.QueryRow(ctx, p.q.GetAvailableSpace, userId).Scan(&aSpace)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, service.ErrNotFound
		}
		return -1, ToServiceError(err)
	}
	return aSpace, nil
}

func (p *PostgresRepo) DelMedia(ctx context.Context, userId, mediaId uuid.UUID) error {
	t, err := p.pool.Exec(ctx, p.q.DelData, userId, mediaId)
	if err != nil {
		return ToServiceError(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFound
	}
	return nil
}

func (p *PostgresRepo) NewProfile(ctx context.Context, userId uuid.UUID) error {
	t, err := p.pool.Exec(ctx, p.q.NewProfile, userId)
	if err != nil {
		return ToServiceError(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrAlreadyExists
	}
	return nil
}

func (p *PostgresRepo) GetProfileInfo(ctx context.Context, userId uuid.UUID) (*service.ProfileInfo, error) {
	var info service.ProfileInfo
	err := p.pool.QueryRow(ctx, p.q.GetProfileInfo, userId).Scan(
		&info.MaxSpace, &info.SpaceFilled, &info.FilesSaved,
	)
	if err != nil {
		return nil, ToServiceError(err)
	}
	return &info, nil
}

func (p *PostgresRepo) GetProfileMediaInfo(ctx context.Context, userId, fromId uuid.UUID, quantity int) ([]service.MediaInfo, error) {
	data, err := repo.GetSliceQueryByPos[service.MediaInfo](ctx, p.pool, p.q.GetProfileData, userId, fromId, quantity)
	if err != nil {
		return nil, ToServiceError(err)
	}
	return data, nil
}
