package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DatabaseItf interface {
	Begin() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error

	PreparexContext(ctx context.Context, query string) *sqlx.Stmt

	GetContextStmt(ctx context.Context, stmt *sqlx.Stmt, dest interface{}, args ...interface{}) error
	SelectContextStmt(ctx context.Context, stmt *sqlx.Stmt, dest interface{}, args ...interface{}) error

	ExecContextStmt(ctx context.Context, stmt *sqlx.Stmt, args ...interface{}) error
	ExecContextStmtTx(ctx context.Context, tx *sql.Tx, stmt *sqlx.Stmt, args ...interface{}) error
}
