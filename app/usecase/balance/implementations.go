package usecasebalance

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
)

func (u usecase) ReadBalanceByUsername(ctx context.Context, req ReadBalanceByUsernameRequest) (resp ReadBalanceByUsernameResponse, err error) {
	balance, err := u.balance.GetBalanceByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return resp, err
	}

	return ReadBalanceByUsernameResponse{
		Balance: balance.Amount,
	}, nil
}

func (u usecase) TopupBalance(ctx context.Context, req TopupBalanceRequest) (resp TopupBalanceResponse, err error) {
	if req.Amount < 0 || req.Amount > 10000000 {
		return TopupBalanceResponse{
			Code: http.StatusBadRequest,
		}, fmt.Errorf("invalid topup amount")
	}

	err = u.balance.GrantBalanceByUsername(ctx, domainbalance.Balance{
		Username: req.Username,
		Amount:   req.Amount,
	})
	if err != nil {
		return TopupBalanceResponse{
			Code: http.StatusBadRequest,
		}, err
	}

	return TopupBalanceResponse{
		Code: http.StatusNoContent,
	}, nil
}

func (u usecase) TransferBalance(ctx context.Context, req TransferBalanceRequest) (resp TransferBalanceResponse, err error) {
	balance, err := u.balance.GetBalanceByUsername(ctx, req.Username)
	if err != nil {
		return TransferBalanceResponse{
			Code: http.StatusBadRequest,
		}, nil
	}

	if balance.Amount-req.Amount < 0 {
		return TransferBalanceResponse{
			Code: http.StatusBadRequest,
		}, fmt.Errorf("insufficient balance")
	}

	destination, err := u.auth.GetUserByUsername(ctx, req.ToUsername)
	if err != nil {
		return TransferBalanceResponse{
			Code: http.StatusNotFound,
		}, nil
	}

	fmt.Println(destination)

	return TransferBalanceResponse{
		Code: http.StatusNoContent,
	}, nil
}
