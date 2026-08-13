package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotFound            = errors.New("not found")
	ErrChatNotFound        = errors.New("chat not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmojiNotFound       = errors.New("emoji not found")
	ErrChatNotFoundOrEmpty = errors.New("chat not found or empty")
	ErrNotFoundOrForbidden = errors.New("forbidden or not found")
	ErrForbidden           = errors.New("forbidden")
	ErrInternal            = errors.New("internal")
)

func getErrorMsg(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "messages_chat_id_fkey":
				return ErrChatNotFound
			case "messages_sender_id_fkey":
				return ErrUserNotFound
			case "chatmembers_user_id_fkey":
				return ErrUserNotFound
			case "message_reactions_emoji_fkey":
				return ErrEmojiNotFound
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
