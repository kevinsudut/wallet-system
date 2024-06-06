package domainbalance

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
	getBalanceByUserId               *sqlx.Stmt
	getLatestHistoryByUserId         *sqlx.Stmt
	getHistorySummaryByUserIdAndType *sqlx.Stmt
	grantBalanceByUserId             *sqlx.Stmt
	deductBalanceByUserId            *sqlx.Stmt
	insertHistory                    *sqlx.Stmt
	updateHistorySummaryById         *sqlx.Stmt
}

func Init(db database.DatabaseItf, redis redis.RedisItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		redis: redis,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			getBalanceByUserId:               db.PreparexContext(ctx, queryGetBalanceByUserId),
			getLatestHistoryByUserId:         db.PreparexContext(ctx, queryGetLatestHistoryByUserId),
			getHistorySummaryByUserIdAndType: db.PreparexContext(ctx, queryGetHistorySummaryByUserIdAndType),
			grantBalanceByUserId:             db.PreparexContext(ctx, queryGrantBalanceByUserId),
			deductBalanceByUserId:            db.PreparexContext(ctx, queryDeductBalanceByUserId),
			insertHistory:                    db.PreparexContext(ctx, queryInsertHistory),
			updateHistorySummaryById:         db.PreparexContext(ctx, queryUpdateHistorySummaryById),
		},
		singleflight: singleflight.Init(),
	}
}
