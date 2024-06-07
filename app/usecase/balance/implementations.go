package usecasebalance

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	domainbalance "github.com/kevinsudut/wallet-system/app/domain/balance"
	"github.com/kevinsudut/wallet-system/app/entity"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
)

func (u usecase) ReadBalanceByUserId(ctx context.Context, req ReadBalanceByUserIdRequest) (resp ReadBalanceByUserIdResponse, err error) {
	balance, err := u.balance.GetBalanceByUserId(ctx, req.UserId)
	if err != nil && err != sql.ErrNoRows {
		log.Errorln("ReadBalanceByUserId.GetBalanceByUserId", err)
		return ReadBalanceByUserIdResponse{
			Code: http.StatusBadRequest,
		}, err
	}

	return ReadBalanceByUserIdResponse{
		Code:    http.StatusOK,
		Balance: balance.Amount,
	}, nil
}

func (u usecase) TopupBalance(ctx context.Context, req TopupBalanceRequest) (resp TopupBalanceResponse, err error) {
	if req.Amount < 0 || req.Amount > 10000000 {
		return TopupBalanceResponse{
			Code: http.StatusBadRequest,
		}, fmt.Errorf("invalid topup amount")
	}

	err = u.balance.GrantBalanceByUserId(ctx, entity.Balance{
		UserId: req.UserId,
		Amount: req.Amount,
	})
	if err != nil {
		log.Errorln("TopupBalance.GrantBalanceByUserId", err)
		return TopupBalanceResponse{
			Code: http.StatusBadRequest,
		}, err
	}

	return TopupBalanceResponse{
		Code: http.StatusNoContent,
	}, nil
}

func (u usecase) TransferBalance(ctx context.Context, req TransferBalanceRequest) (resp TransferBalanceResponse, err error) {
	balance, err := u.balance.GetBalanceByUserId(ctx, req.UserId)
	if err != nil {
		log.Errorln("TransferBalance.GetBalanceByUserId", err)
		return TransferBalanceResponse{
			Code: http.StatusBadRequest,
		}, err
	}

	if balance.Amount-req.Amount < 0 {
		return TransferBalanceResponse{
			Code: http.StatusBadRequest,
		}, fmt.Errorf("insufficient balance")
	}

	toUser, err := u.auth.GetUserByUsername(ctx, req.ToUsername)
	if err != nil {
		log.Errorln("TransferBalance.GetUserByUsername", err)
		return TransferBalanceResponse{
			Code: http.StatusNotFound,
		}, err
	}

	err = u.balance.DisburmentBalance(ctx, domainbalance.DisburmentBalanceRequest{
		UserId:   req.UserId,
		ToUserId: toUser.Id,
		Amount:   req.Amount,
	})
	if err != nil {
		log.Errorln("TransferBalance.DisburmentBalance", err)
		return TransferBalanceResponse{
			Code: http.StatusBadRequest,
		}, err
	}

	return TransferBalanceResponse{
		Code: http.StatusNoContent,
	}, nil
}
