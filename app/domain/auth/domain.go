package domainauth

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
)

type domain struct {
	db    database.DatabaseItf
	cache lrucache.LRUCacheItf
	stmts databaseStmts
}

type databaseStmts struct {
	getUserByUsername *sqlx.Stmt
	insertUser        *sqlx.Stmt
}

func Init(db database.DatabaseItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			getUserByUsername: db.PreparexContext(ctx, queryGetUserByUsername),
			insertUser:        db.PreparexContext(ctx, queryInsertUser),
		},
	}
}
