package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DaiYuANg/arcgo/dbx"
	sqlitedialect "github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
	_ "modernc.org/sqlite"
)

func OpenSQLite(dsn string, opts ...dbx.Option) (*dbx.DB, error) {
	if dsn == "" {
		dsn = "file:backend?mode=memory&cache=shared"
	}
	db, err := dbx.Open(
		dbx.WithDriver("sqlite"),
		dbx.WithDSN(dsn),
		dbx.WithDialect(sqlitedialect.New()),
		dbx.ApplyOptions(opts...),
	)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.ExecContext(context.Background(), `PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func DefaultOpts(logger *slog.Logger) []dbx.Option {
	if logger == nil {
		return nil
	}
	return []dbx.Option{
		dbx.WithLogger(logger),
		dbx.WithDebug(false),
	}
}
