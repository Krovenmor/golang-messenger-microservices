package migrations

import (
	"MyMessenger/pkg/repo"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var EmbedMigrations embed.FS

func MakeMigrations(pool *pgxpool.Pool) error {
	return repo.MakeMigrations(pool, &EmbedMigrations)
}
