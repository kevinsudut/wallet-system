package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (db database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.conn.GetContext(ctx, dest, query, args...)
}

func (db database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.conn.SelectContext(ctx, dest, query, args...)
}

func (db database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.conn.ExecContext(ctx, query, args...)
}

func (db database) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

func (db database) ExecContextTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(ctx, query, args...)
}

func (db database) PreparexContext(ctx context.Context, query string) *sqlx.Stmt {
	stmt, err := db.conn.PreparexContext(ctx, query)
	if err != nil {
		log.Panicln("invalid prepare query", query)
	}
	return stmt
}
