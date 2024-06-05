package usecasebalance

import "context"

type UsecaseItf interface {
	ReadBalanceByUserId(ctx context.Context, req ReadBalanceByUserIdRequest) (resp ReadBalanceByUserIdResponse, err error)
	TopupBalance(ctx context.Context, req TopupBalanceRequest) (resp TopupBalanceResponse, err error)
	TransferBalance(ctx context.Context, req TransferBalanceRequest) (resp TransferBalanceResponse, err error)
}
