package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (db database) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

func (db database) PreparexContext(ctx context.Context, query string) *sqlx.Stmt {
	stmt, err := db.conn.PreparexContext(ctx, query)
	if err != nil {
		log.Panicln("invalid prepare query", query)
	}
	return stmt
}

func (db database) GetContextStmt(ctx context.Context, stmt *sqlx.Stmt, dest interface{}, args ...interface{}) error {
	return stmt.GetContext(ctx, dest, args...)
}

func (db database) SelectContextStmt(ctx context.Context, stmt *sqlx.Stmt, dest interface{}, args ...interface{}) error {
	return stmt.SelectContext(ctx, dest, args...)
}

func (db database) ExecContextStmt(ctx context.Context, stmt *sqlx.Stmt, args ...interface{}) error {
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed ExecContextStmt %s", err.Error())
	}

	return nil
}

func (db database) ExecContextStmtTx(ctx context.Context, tx *sql.Tx, stmt *sqlx.Stmt, args ...interface{}) error {
	result, err := tx.StmtContext(ctx, stmt.Stmt).ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed ExecContextStmtTx %s", err.Error())
	}

	return nil
}
