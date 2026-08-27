package postgres

import (
	"MyMessenger/services/media/internal/service"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

func ToServiceError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return service.ErrNotFound
	}
	log.Printf("Repo, ToServiceError: Unknown err: %q", err)
	return service.ErrUnknown
}
