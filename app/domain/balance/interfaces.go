package domainbalance

import "context"

type DomainItf interface {
	GetBalanceByUserId(ctx context.Context, userId string) (resp Balance, err error)
	GrantBalanceByUserId(ctx context.Context, balance Balance) (err error)
	DisburmentBalance(ctx context.Context, req DisburmentBalanceRequest) (err error)

	GetLatestHistoryByUserId(ctx context.Context, userId string) (resp []History, err error)
	GetHistorySummaryByUserIdAndType(ctx context.Context, userId string, historyType int) (resp []HistorySummary, err error)
}
