package domainbalance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (d domain) GetBalanceByUserId(ctx context.Context, userId string) (resp Balance, err error) {
	defer func() {
		if err == nil && resp.UserId == "" {
			err = sql.ErrNoRows
		}
	}()

	balance, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetBalanceByUserId, userId), func() (interface{}, error) {
		var resp Balance
		balance, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetBalanceByUserId, userId), time.Minute*5, func() (interface{}, error) {
			var balance Balance
			err := d.db.GetContextStmt(ctx, d.stmts.getBalanceByUserId, &balance, userId)
			if err != nil && err != sql.ErrNoRows {
				return balance, err
			}

			return balance, nil
		})
		if err != nil {
			return resp, err
		}

		return balance.Value().(Balance), nil
	})
	if err != nil {
		return resp, err
	}

	return balance.(Balance), nil
}

func (d domain) grantBalanceByUserId(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.grantBalanceByUserId, balance.UserId, balance.Amount)
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) deductBalanceByUserId(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.deductBalanceByUserId, balance.Amount, balance.UserId)
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) insertHistory(ctx context.Context, tx *sql.Tx, history History) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.insertHistory, history.Id, history.UserId, history.TargetUserId, history.Amount, history.Type, history.Notes)
	if err != nil {
		return err
	}

	err = d.updateHistorySummary(ctx, tx, HistorySummary{
		UserId:       history.UserId,
		TargetUserId: history.TargetUserId,
		Amount:       history.Amount,
		Type:         history.Type,
	})
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, history.UserId))

	return nil
}

func (d domain) updateHistorySummary(ctx context.Context, tx *sql.Tx, historySummary HistorySummary) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.updateHistorySummaryById, historySummary.GetId(), historySummary.UserId, historySummary.TargetUserId, historySummary.Amount, historySummary.Type)
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, historySummary.UserId, historySummary.Type))

	return nil
}

func (d domain) GrantBalanceByUserId(ctx context.Context, balance Balance) (err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = d.grantBalanceByUserId(ctx, tx, balance)
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:           uuid.NewString(),
		UserId:       balance.UserId,
		TargetUserId: balance.UserId,
		Amount:       balance.Amount,
		Type:         int(CREDIT),
		Notes:        "Top-up money",
	})
	if err != nil {
		return err
	}

	return nil
}

func (d domain) DisburmentBalance(ctx context.Context, req DisburmentBalanceRequest) (err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = d.grantBalanceByUserId(ctx, tx, Balance{
		UserId: req.ToUserId,
		Amount: req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.deductBalanceByUserId(ctx, tx, Balance{
		UserId: req.UserId,
		Amount: req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:           uuid.NewString(),
		UserId:       req.ToUserId,
		TargetUserId: req.UserId,
		Amount:       req.Amount,
		Type:         int(CREDIT),
		Notes:        fmt.Sprintf("Receive money from %s", req.UserId),
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:           uuid.NewString(),
		UserId:       req.UserId,
		TargetUserId: req.ToUserId,
		Amount:       req.Amount,
		Type:         int(DEBIT),
		Notes:        fmt.Sprintf("Transfer money to %s", req.ToUserId),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d domain) GetLatestHistoryByUserId(ctx context.Context, userId string) (resp []History, err error) {
	defer func() {
		if err == nil {
			for i := range resp {
				resp[i].NormalizeAmount()
			}
		}
	}()

	histories, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetLatestHistoryByUserId, userId), func() (interface{}, error) {
		var resp []History
		histories, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, userId), time.Minute*5, func() (interface{}, error) {
			var history []History
			err := d.db.SelectContextStmt(ctx, d.stmts.getLatestHistoryByUserId, &history, userId)
			if err != nil {
				return history, err
			}

			return history, nil
		})
		if err != nil {
			return resp, err
		}

		return histories.Value().([]History), nil
	})
	if err != nil {
		return resp, err
	}

	return histories.([]History), nil
}

func (d domain) GetHistorySummaryByUserIdAndType(ctx context.Context, userId string, historyType int) (resp []HistorySummary, err error) {
	historySummaries, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetHistorySummaryByUserIdAndType, userId, historyType), func() (interface{}, error) {
		var resp []HistorySummary
		historySummaries, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, userId, historyType), time.Minute*5, func() (interface{}, error) {
			var historySummary []HistorySummary
			err := d.db.SelectContextStmt(ctx, d.stmts.getHistorySummaryByUserIdAndType, &historySummary, userId, historyType)
			if err != nil {
				return historySummary, err
			}

			return historySummary, nil
		})
		if err != nil {
			return resp, err
		}

		return historySummaries.Value().([]HistorySummary), nil
	})
	if err != nil {
		return resp, err
	}

	return historySummaries.([]HistorySummary), nil
}
