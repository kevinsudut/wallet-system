package domainbalance

import "context"

type DomainItf interface {
	GetBalanceByUsername(ctx context.Context, username string) (resp Balance, err error)
	GrantBalanceByUsername(ctx context.Context, balance Balance) (err error)
	DisburmentBalance(ctx context.Context, req DisburmentBalanceRequest) (err error)
}
