package repo

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func MakeMigrations(pool *pgxpool.Pool, migrations *embed.FS) error {
	if pool == nil {
		return fmt.Errorf("MakeMigrations(): pool == nil")
	}
	if migrations == nil {
		return fmt.Errorf("MakeMigrations(): migrations == nil")
	}

	db := stdlib.OpenDBFromPool(pool)
	if db == nil {
		return fmt.Errorf("MakeMigrations(): trouble with OpenDBFromPool")
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Trouble with db.Close(), err: %q", err)
		}
	}()

	goose.SetBaseFS(*migrations)
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
