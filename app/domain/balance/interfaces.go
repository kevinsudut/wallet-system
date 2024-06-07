package domainbalance

import (
	"context"

	"github.com/kevinsudut/wallet-system/app/entity"
)

type DomainItf interface {
	GetBalanceByUserId(ctx context.Context, userId string) (resp entity.Balance, err error)
	GrantBalanceByUserId(ctx context.Context, balance entity.Balance) (err error)
	DisburmentBalance(ctx context.Context, req DisburmentBalanceRequest) (err error)

	GetLatestHistoryByUserId(ctx context.Context, userId string) (resp []entity.History, err error)
	GetHistorySummaryByUserIdAndType(ctx context.Context, userId string, historyType int) (resp []entity.HistorySummary, err error)
}
