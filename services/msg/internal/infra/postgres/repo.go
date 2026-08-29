package postgres

import (
	"MyMessenger/pkg/repo"
	"MyMessenger/services/msg/internal/infra/postgres/queries"
	"MyMessenger/services/msg/internal/service"
	"context"
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

func (r *PostagreRepo) NewProfile(ctx context.Context, userId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.NewProfile, userId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) NewChat(ctx context.Context, chatId, fUser, sUser uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.NewChat, chatId, fUser, sUser)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) NewMessage(ctx context.Context, chatId uuid.UUID, msg *service.Message) error {
	t, err := r.pool.Exec(ctx, r.q.PostMessage, msg.MessageId, chatId, msg.SenderId, msg.Message, msg.SenderKey, msg.ReceiverKey, msg.Nonce, msg.ReplyToId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrAlreadyExists
	}
	return nil
}

func (r *PostagreRepo) getReactions(ctx context.Context, msgIds []uuid.UUID) ([]Reaction, error) {
	reactions, err := repo.GetSliceQueryByPos[Reaction](ctx, r.pool, r.q.GetReactions, msgIds)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return reactions, nil
}

func (r *PostagreRepo) GetMessage(ctx context.Context, chatId, msgId uuid.UUID) (*service.Message, error) {
	msg, err := repo.GetQueryByPos[Message](ctx, r.pool, r.q.GetMessage, chatId, msgId)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	reactions, err := r.getReactions(ctx, []uuid.UUID{msg.MessageId})
	if err != nil {
		return nil, getErrorMsg(err)
	}
	sReactions := ToServiceReactions(reactions)
	sMessage := ToServiceMessage(msg, sReactions)
	return &sMessage, nil
}

func (r *PostagreRepo) RedactMessage(ctx context.Context, chatId, msgId uuid.UUID, msg *service.ToPostMessage) error {
	t, err := r.pool.Exec(ctx, r.q.UpdateMessage,
		chatId, msgId, msg.UserId, msg.Message, msg.ReceiverKey, msg.SenderKey, msg.Nonce,
	)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFoundOrForbidden
	}
	return nil
}

func (r *PostagreRepo) DelMessage(ctx context.Context, chatId, msgId, userId uuid.UUID) error {
	t, err := r.pool.Exec(ctx, r.q.DeleteMessage, chatId, msgId, userId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrNotFoundOrForbidden
	}
	return nil
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
		&info.CreatedAt, &info.Members,
	)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return &info, nil
}

func (r *PostagreRepo) GetChatHistory(ctx context.Context, chatId uuid.UUID, fromId uuid.UUID, q int) ([]service.Message, error) {
	messages, err := repo.GetSliceQueryByPos[Message](ctx, r.pool, r.q.GetChatHistory, chatId, fromId, q)
	if err != nil {
		return nil, getErrorMsg(err)
	}

	sMessages := make([]service.Message, len(messages))
	ids := make([]uuid.UUID, len(messages))
	for i := range messages {
		ids[i] = messages[i].MessageId
		sMessages[i] = ToServiceMessage(&messages[i], []service.Reaction{})
	}

	reactions, err := r.getReactions(ctx, ids)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	if len(reactions) == 0 {
		return sMessages, nil
	}

	reactionIdx := 0
	numReactions := len(reactions)

	for i := range sMessages {
		for reactionIdx < numReactions && reactions[reactionIdx].MessageId == sMessages[i].MessageId {
			sMessages[i].Reactions = append(sMessages[i].Reactions, ToServiceReaction(reactions[reactionIdx]))
			reactionIdx++
		}
	}

	return sMessages, nil
}

func (r *PostagreRepo) GetEmojis(ctx context.Context) ([]string, error) {
	sl, err := repo.GetSliceQueryByType[string](ctx, r.pool, r.q.GetEmojis)
	if err != nil {
		return nil, getErrorMsg(err)
	}
	return sl, nil
}

func (r *PostagreRepo) NewReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error {
	t, err := r.pool.Exec(ctx, r.q.NewReaction, msgId, userId, emoji, chatId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrForbidden
	}
	return nil
}

func (r *PostagreRepo) DelReaction(ctx context.Context, userId, chatId, msgId uuid.UUID, emoji string) error {
	t, err := r.pool.Exec(ctx, r.q.DelReaction, msgId, userId, emoji, chatId)
	if err != nil {
		return getErrorMsg(err)
	}
	if t.RowsAffected() == 0 {
		return service.ErrForbidden
	}
	return nil
}

func (r *PostagreRepo) IsProfileInChat(ctx context.Context, userId, chatId uuid.UUID) error {
	var isInChat int
	err := r.pool.QueryRow(ctx, r.q.IsProfileInChat, chatId, userId).Scan(&isInChat)
	if err != nil {
		return getErrorMsg(err)
	}
	if isInChat != 1 {
		return service.ErrNotFound
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
