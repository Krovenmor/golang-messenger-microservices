package repo

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type gooseLogger struct {
	*zap.SugaredLogger
}

func (l gooseLogger) Printf(format string, v ...any) {
	l.Infof(strings.TrimSuffix(format, "\n"), v...)
}

func MakeMigrations(pool *pgxpool.Pool, migrations *embed.FS) (err error) {
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
		err = errors.Join(err, db.Close())
	}()

	goose.SetLogger(gooseLogger{zap.S()})
	goose.SetBaseFS(*migrations)
	err = goose.SetDialect("postgres")
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
