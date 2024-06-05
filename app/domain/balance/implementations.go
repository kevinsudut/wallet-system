package domainbalance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func (d domain) GetBalanceByUserId(ctx context.Context, userId string) (resp Balance, err error) {
	defer func() {
		if err == nil && resp.UserId == "" {
			err = sql.ErrNoRows
		}
	}()

	balance, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetBalanceByUserId, userId), func() (interface{}, error) {
		var resp Balance
		balance, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetBalanceByUserId, userId), time.Minute*5, func() (string, error) {
			var balance Balance
			err := d.stmts.getBalanceByUserId.GetContext(ctx, &balance, userId)
			if err != nil && err != sql.ErrNoRows {
				return "", err
			}

			return jsoniter.MarshalToString(balance)
		})
		if err != nil {
			return resp, err
		}

		err = jsoniter.UnmarshalFromString(balance.Value(), &resp)
		if err != nil {
			return resp, err
		}

		return resp, nil
	})
	if err != nil {
		return resp, err
	}

	return balance.(Balance), nil
}

func (d domain) grantBalanceByUserId(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	result, err := tx.ExecContext(ctx, queryGrantBalanceByUserId, balance.UserId, balance.Amount)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) deductBalanceByUserId(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	result, err := tx.ExecContext(ctx, queryDeductBalanceByUserId, balance.Amount, balance.UserId)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) insertHistory(ctx context.Context, tx *sql.Tx, history History) (err error) {
	result, err := tx.ExecContext(ctx, queryInsertHistory, history.Id, history.UserId, history.TargetUserId, history.Amount, history.Type, history.Notes)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
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

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetLatestHistoryByUserId, history.UserId))

	return nil
}

func (d domain) updateHistorySummary(ctx context.Context, tx *sql.Tx, historySummary HistorySummary) (err error) {
	result, err := tx.ExecContext(ctx, queryUpdateHistorySummaryById, historySummary.GetId(), historySummary.UserId, historySummary.TargetUserId, historySummary.Amount, historySummary.Type)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetHistorySummaryByUserIdAndType, historySummary.UserId, historySummary.Type))

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
		UserId: req.ToUserId,
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
		TargetUserId: req.UserId,
		Amount:       req.Amount,
		Type:         int(DEBIT),
		Notes:        fmt.Sprintf("Transfer money to %s", req.UserId),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d domain) GetLatestHistoryByUserId(ctx context.Context, userId string) (resp []History, err error) {
	defer func() {
		if err != nil {
			for i := range resp {
				resp[i].NormalizeAmount()
			}
		}
	}()

	histories, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetLatestHistoryByUserId, userId), func() (interface{}, error) {
		var resp []History
		histories, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetLatestHistoryByUserId, userId), time.Minute*5, func() (string, error) {
			var history []History
			err := d.stmts.getLatestHistoryByUserId.SelectContext(ctx, &history, userId)
			if err != nil {
				return "", err
			}

			return jsoniter.MarshalToString(history)
		})
		if err != nil {
			return resp, err
		}

		err = jsoniter.UnmarshalFromString(histories.Value(), &resp)
		if err != nil {
			return resp, err
		}

		return resp, nil
	})
	if err != nil {
		return resp, err
	}

	return histories.([]History), nil
}

func (d domain) GetHistorySummaryByUserIdAndType(ctx context.Context, userId string, historyType int) (resp []HistorySummary, err error) {
	historySummaries, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetHistorySummaryByUserIdAndType, userId), func() (interface{}, error) {
		var resp []HistorySummary
		historySummaries, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetHistorySummaryByUserIdAndType, userId, historyType), time.Minute*5, func() (string, error) {
			var historySummary []HistorySummary
			err := d.stmts.getHistorySummaryByUserIdAndType.SelectContext(ctx, &historySummary, userId, historyType)
			if err != nil {
				return "", err
			}

			return jsoniter.MarshalToString(historySummary)
		})
		if err != nil {
			return resp, err
		}

		err = jsoniter.UnmarshalFromString(historySummaries.Value(), &resp)
		if err != nil {
			return resp, err
		}

		return resp, nil
	})
	if err != nil {
		return resp, err
	}

	return historySummaries.([]HistorySummary), nil
}
