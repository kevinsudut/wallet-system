package usecasebalance

import "context"

type UsecaseItf interface {
	ReadBalanceByUsername(ctx context.Context, req ReadBalanceByUsernameRequest) (resp ReadBalanceByUsernameResponse, err error)
	TopupBalance(ctx context.Context, req TopupBalanceRequest) (resp TopupBalanceResponse, err error)
	TransferBalance(ctx context.Context, req TransferBalanceRequest) (resp TransferBalanceResponse, err error)
}
