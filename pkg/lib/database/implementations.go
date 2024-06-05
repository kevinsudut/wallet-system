package database

import (
	"context"
	"database/sql"
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
