package infra

import (
	"MyMessenger/services/auth/internal/infra/migrations"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func MakeMigrations(pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("MakeMigrations(): pool == nil")
	}

	db := stdlib.OpenDBFromPool(pool)
	if db == nil {
		return fmt.Errorf("MakeMigrations(): trouble with OpenDBFromPool")
	}
	defer db.Close()

	goose.SetBaseFS(migrations.EmbedMigrations)
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = goose.UpContext(ctx, db, ".")
	if err != nil {
		return err
	}

	return nil
}
