package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DatabaseItf interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	Begin() (*sql.Tx, error)
	ExecContextTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)

	PreparexContext(ctx context.Context, query string) *sqlx.Stmt
}
