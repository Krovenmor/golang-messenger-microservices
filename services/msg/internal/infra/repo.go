package infra

import (
	"MyMessenger/pkg/repo"
	"MyMessenger/services/msg/internal/infra/queries"
	"MyMessenger/services/msg/internal/service"
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
	)
	if err != nil {
		return getErrorMsg(err)
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

func (r *PostagreRepo) NewMessage(ctx context.Context, chatId uuid.UUID, msg *service.Message) error {
	t, err := r.pool.Exec(ctx, r.q.PostMessage, msg.MessageId, chatId, msg.SenderId, msg.Message)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*service.Message, error) {
	var msg service.Message
	err := r.pool.QueryRow(ctx, r.q.GetMessage, chatId, msgId).Scan(
		&msg.MessageId, &msg.SenderId, &msg.Message, &msg.CreatedAt,
		&msg.IsRedacted, &msg.IsDeleted, &msg.RedactedAt,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return &msg, nil
}

func (r *PostagreRepo) RedactMessage(ctx context.Context, chatId, msgId, userId uuid.UUID, newText string) error {
	t, err := r.pool.Exec(ctx, r.q.UpdateMessage, chatId, msgId, userId, newText)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}

func (r *PostagreRepo) DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.DeleteMessage, chatId, msgId, userId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}

func (r *PostagreRepo) GetProfileById(ctx context.Context, userId uuid.UUID) (*service.Profile, error) {
	return getProfileByVal(ctx, r.pool, r.q.GetProfileId, userId)
}

func (r *PostagreRepo) GetProfileByUserName(ctx context.Context, username string) (*service.Profile, error) {
	return getProfileByVal(ctx, r.pool, r.q.GetProfileUserName, username)
}

func (r *PostagreRepo) GetChats(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	collectFunc := func(row pgx.CollectableRow) (uuid.UUID, error) {
		var id uuid.UUID
		err := row.Scan(&id)
		return id, err
	}
	chats, err := repo.GetSliceQueryByFunc(ctx, r.pool, r.q.GetChats, collectFunc, userId)
	if err != nil {
		return chats, getErrorMsg(err)
	}
	return chats, nil
}

func (r *PostagreRepo) GetChatsExtended(ctx context.Context, userId uuid.UUID) ([]service.ChatFullInfo, error) {
	chats, err := repo.GetSliceQueryByPos[service.ChatFullInfo](ctx, r.pool, r.q.GetChatsExtended, userId)
	if err != nil {
		return chats, getErrorMsg(err)
	}
	return chats, nil
}

func (r *PostagreRepo) GetChatMembers(ctx context.Context, chatId uuid.UUID) ([]service.ChatMember, error) {
	members, err := repo.GetSliceQueryByPos[service.ChatMember](ctx, r.pool, r.q.GetChatMembers, chatId)
	if err != nil {
		return members, getErrorMsg(err)
	}
	return members, nil
}

func (r *PostagreRepo) GetChatInfo(ctx context.Context, chatId uuid.UUID) (*service.ChatInfo, error) {
	var info service.ChatInfo
	err := r.pool.QueryRow(ctx, r.q.GetChatInfo, chatId).Scan(
		&info.CreatedAt,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	info.ChatMembers, err = r.GetChatMembers(ctx, chatId)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *PostagreRepo) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]service.Message, error) {
	messages, err := repo.GetSliceQueryByPos[service.Message](ctx, r.pool, r.q.GetChatHistory, chatId, fromId, q)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return []service.Message{}, nil
		}
		return messages, getErrorMsg(err)
	}
	return messages, nil
}

func (r *PostagreRepo) IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error {
	var isInChat int
	err := r.pool.QueryRow(ctx, r.q.IsProfileInChat, chatId, userId).Scan(&isInChat)
	if err != nil {
		return getErrorMsg(err)
	}
	if isInChat != 1 {
		return ErrNotFound
	}
	return nil
}

func (r *PostagreRepo) IsProfilesHaveAPrivateChat(ctx context.Context, userIdF, userIdS uuid.UUID) (uuid.UUID, error) {
	var chatId uuid.UUID
	err := r.pool.QueryRow(ctx, r.q.GetPrivateChatBetweenTwoPeoples, userIdF, userIdS).Scan(
		&chatId,
	)
	if err != nil {
		return chatId, getErrorMsg(err)
	}
	return chatId, nil
}
