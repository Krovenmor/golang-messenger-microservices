package postgres

import (
	"MyMessenger/services/auth/internal/domain"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	uniqueViolationCode = "23505"

	loginUniqueConstraint = "users_login_key"
	emailUniqueConstraint = "idx_unique_email"
)

func betterError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return uniqueViolationError(pgErr.ConstraintName)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func uniqueViolationError(constraint string) error {
	switch constraint {
	case loginUniqueConstraint:
		return domain.ErrLoginTaken
	case emailUniqueConstraint:
		return domain.ErrEmailTaken
	default:
		return domain.ErrAlreadyExists
	}
}
