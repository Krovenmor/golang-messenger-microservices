package postgres

import (
	"MyMessenger/services/profile/internal/service"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func getErrorMsg(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return service.ErrAlreadyExists
		}
	}
	log.Printf("Not defined err: %q", err)
	return service.ErrUnknown
}
