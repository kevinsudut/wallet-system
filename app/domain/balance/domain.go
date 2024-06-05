package domainbalance

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
	getBalanceByUsername               *sqlx.Stmt
	getLatestHistoryByUsername         *sqlx.Stmt
	getHistorySummaryByUsernameAndType *sqlx.Stmt
}

func Init(db database.DatabaseItf) DomainItf {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return &domain{
		db:    db,
		cache: lrucache.Init(),
		stmts: databaseStmts{
			getBalanceByUsername:               db.PreparexContext(ctx, queryGetBalanceByUsername),
			getLatestHistoryByUsername:         db.PreparexContext(ctx, queryGetLatestHistoryByUsername),
			getHistorySummaryByUsernameAndType: db.PreparexContext(ctx, queryGetHistorySummaryByUsernameAndType),
		},
	}
}
