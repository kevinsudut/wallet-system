package domainauth

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kevinsudut/wallet-system/pkg/helper/singleflight"
	"github.com/kevinsudut/wallet-system/pkg/lib/database"
	lrucache "github.com/kevinsudut/wallet-system/pkg/lib/lru-cache"
	"github.com/kevinsudut/wallet-system/pkg/lib/redis"
)

type domain struct {
	db           database.DatabaseItf
	redis        redis.RedisItf
	cache        lrucache.LRUCacheItf
	stmts        databaseStmts
	singleflight singleflight.SingleFlightItf
}

type databaseStmts struct {
	insertUser        *sqlx.Stmt
	getUserById       *sqlx.Stmt
	getUserByUsername *sqlx.Stmt
}

func Init(db database.DatabaseItf, redis redis.RedisItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		redis: redis,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			insertUser:        db.PreparexContext(ctx, queryInsertUser),
			getUserById:       db.PreparexContext(ctx, queryGetUserById),
			getUserByUsername: db.PreparexContext(ctx, queryGetUserByUsername),
		},
		singleflight: singleflight.Init(),
	}
}
