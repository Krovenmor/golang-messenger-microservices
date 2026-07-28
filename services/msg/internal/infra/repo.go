package infra

import (
	"MyMessenger/services/msg/internal/infra/queries"
	"MyMessenger/services/msg/internal/service"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotFound            = errors.New("not found")
	ErrChatNotFound        = errors.New("chat not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrChatNotFoundOrEmpty = errors.New("chat not found or empty")
)

func getErrorMsg(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "messages_chatid_fkey":
				return ErrChatNotFound
			case "messages_sentid_fkey":
				return ErrUserNotFound
			case "chatmembers_userid_fkey":
				return ErrUserNotFound
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

func (r *PostagreRepo) NewProfile(ctx context.Context, profile service.Profile) error {
	t, err := r.pool.Exec(ctx, r.q.NewProfile, profile.UserId, profile.Name, profile.PublicKey, profile.PrivateKey, profile.KDFSalt)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.NewChat, chatId, fUser, sUser)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) PostMessage(ctx context.Context, chatId uuid.UUID, msg service.Message) error {
	t, err := r.pool.Exec(ctx, r.q.PostMessage, msg.MessageId, chatId, msg.SenderId, msg.Message)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) GetProfile(ctx context.Context, userId uuid.UUID) (*service.Profile, error) {
	var profile service.Profile
	err := r.pool.QueryRow(ctx, r.q.GetProfile, userId).Scan(
		&profile.UserId,
		&profile.Name,
		&profile.PublicKey,
		&profile.PrivateKey,
		&profile.KDFSalt,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return &profile, nil
}

func (r *PostagreRepo) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]service.Message, error) {
	messages := []service.Message{}
	rows, err := r.pool.Query(ctx, r.q.GetChatHistory, chatId, fromId, q)
	if err != nil {
		return messages, getErrorMsg(err)
	}
	defer rows.Close()

	messages, err = pgx.CollectRows(rows, pgx.RowToStructByName[service.Message])
	if err != nil {
		return nil, getErrorMsg(err)
	}

	if len(messages) == 0 {
		return nil, ErrChatNotFoundOrEmpty
	}

	return messages, nil
}

func (r *PostagreRepo) IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.IsProfileInChat, chatId, userId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
