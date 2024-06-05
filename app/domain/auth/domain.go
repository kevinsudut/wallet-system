package domainauth

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
)

type domain struct {
	db           database.DatabaseItf
	cache        lrucache.LRUCacheItf
	stmts        databaseStmts
	singleflight singleflight.SingleFlightItf
}

type databaseStmts struct {
	insertUser        *sqlx.Stmt
	getUserById       *sqlx.Stmt
	getUserByUsername *sqlx.Stmt
}

func Init(db database.DatabaseItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			insertUser:        db.PreparexContext(ctx, queryInsertUser),
			getUserById:       db.PreparexContext(ctx, queryGetUserById),
			getUserByUsername: db.PreparexContext(ctx, queryGetUserByUsername),
		},
		singleflight: singleflight.Init(),
	}
}
