package postgres

import (
	"errors"
	"log"

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
	ErrTooMuch             = errors.New("too much")

	ErrUnknown = errors.New("unknown")
)

func getErrorMsg(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
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

		case "23505":
			return ErrAlreadyExists

		case "23514":
			switch pgErr.ConstraintName {
			case "check_avatars_in_additional":
				return ErrTooMuch
			}
		}

	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	log.Printf("Not defined err: %q", err)
	return ErrUnknown
}
