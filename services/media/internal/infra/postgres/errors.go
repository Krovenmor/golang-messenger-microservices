package postgres

import (
	"MyMessenger/services/media/internal/service"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func ToServiceError(err error) error {
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
	log.Printf("Repo, ToServiceError: Unknown err: %q", err)
	return service.ErrUnknown
}
