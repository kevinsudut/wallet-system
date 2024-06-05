package domainbalance

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
	getBalanceByUserId               *sqlx.Stmt
	getLatestHistoryByUserId         *sqlx.Stmt
	getHistorySummaryByUserIdAndType *sqlx.Stmt
}

func Init(db database.DatabaseItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			getBalanceByUserId:               db.PreparexContext(ctx, queryGetBalanceByUserId),
			getLatestHistoryByUserId:         db.PreparexContext(ctx, queryGetLatestHistoryByUserId),
			getHistorySummaryByUserIdAndType: db.PreparexContext(ctx, queryGetHistorySummaryByUserIdAndType),
		},
		singleflight: singleflight.Init(),
	}
}
