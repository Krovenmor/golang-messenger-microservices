package infra

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
