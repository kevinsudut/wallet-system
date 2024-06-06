package database

import (
	"context"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type database struct {
	conn *sqlx.DB
}

func Init() (DatabaseItf, error) {
	conn, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = conn.PingContext(context.Background())
	if err != nil {
		return nil, err
	}

	return &database{
		conn: conn,
	}, nil
}
