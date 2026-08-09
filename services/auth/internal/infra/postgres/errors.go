package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists     = errors.New("already exists")
	ErrAlreadyExistsUser = errors.New("user already exists")
	ErrNotFound          = errors.New("not found")
)

func betterError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "users_login_key":
				return ErrAlreadyExistsUser
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
