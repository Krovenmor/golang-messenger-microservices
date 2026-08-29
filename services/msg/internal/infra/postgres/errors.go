package postgres

import (
	"MyMessenger/services/msg/internal/service"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func getErrorMsg(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			switch pgErr.ConstraintName {
			case "messages_chat_id_fkey":
				return service.ErrChatNotFound
			case "messages_sender_id_fkey":
				return service.ErrUserNotFound
			case "chatmembers_user_id_fkey":
				return service.ErrUserNotFound
			case "message_reactions_emoji_fkey":
				return service.ErrEmojiNotFound
			}

		case "23505":
			return service.ErrAlreadyExists

		case "23514":
			switch pgErr.ConstraintName {
			case "check_avatars_in_additional":
				return service.ErrTooMuch
			}
		}

	}
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ErrNotFound
	}
	log.Printf("Not defined err: %q", err)
	return service.ErrUnknown
}
