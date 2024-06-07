package domainbalance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/kevinsudut/wallet-system/app/entity"
	"github.com/kevinsudut/wallet-system/app/enum"
)

func (d domain) GetBalanceByUserId(ctx context.Context, userId string) (resp entity.Balance, err error) {
	defer func() {
		if err == nil && resp.UserId == "" {
			err = sql.ErrNoRows
		}
	}()

	balance, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetBalanceByUserId, userId), func() (interface{}, error) {
		var resp entity.Balance
		balance, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetBalanceByUserId, userId), time.Minute*5, func() (interface{}, error) {
			var respRedis entity.Balance
			balanceStr, err := d.redis.Fetch(ctx, fmt.Sprintf(cacheKeyGetBalanceByUserId, userId), time.Duration(time.Minute*30), func() (interface{}, error) {
				var balance entity.Balance
				err := d.db.GetContextStmt(ctx, d.stmts.getBalanceByUserId, &balance, userId)
				if err != nil && err != sql.ErrNoRows {
					return balance, err
				}

				return balance, nil
			})
			if err != nil {
				return respRedis, err
			}

			err = jsoniter.UnmarshalFromString(balanceStr, &respRedis)
			if err != nil {
				return respRedis, err
			}

			return respRedis, nil
		})
		if err != nil {
			return resp, err
		}

		return balance.Value().(entity.Balance), nil
	})
	if err != nil {
		return resp, err
	}

	return balance.(entity.Balance), nil
}

func (d domain) grantBalanceByUserId(ctx context.Context, tx *sql.Tx, balance entity.Balance) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.grantBalanceByUserId, balance.UserId, balance.Amount)
	if err != nil {
		return err
	}

	_, err = d.redis.Delete(ctx, fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) deductBalanceByUserId(ctx context.Context, tx *sql.Tx, balance entity.Balance) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.deductBalanceByUserId, balance.Amount, balance.UserId)
	if err != nil {
		return err
	}

	_, err = d.redis.Delete(ctx, fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetBalanceByUserId, balance.UserId))

	return nil
}

func (d domain) insertHistory(ctx context.Context, tx *sql.Tx, history entity.History) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.insertHistory, history.Id, history.UserId, history.TargetUserId, history.Amount, history.Type, history.Notes)
	if err != nil {
		return err
	}

	err = d.updateHistorySummary(ctx, tx, entity.HistorySummary{
		UserId:       history.UserId,
		TargetUserId: history.TargetUserId,
		Amount:       history.Amount,
		Type:         history.Type,
	})
	if err != nil {
		return err
	}

	_, err = d.redis.Delete(ctx, fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, history.UserId))
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, history.UserId))

	return nil
}

func (d domain) updateHistorySummary(ctx context.Context, tx *sql.Tx, historySummary entity.HistorySummary) (err error) {
	err = d.db.ExecContextStmtTx(ctx, tx, d.stmts.updateHistorySummaryById, historySummary.GetId(), historySummary.UserId, historySummary.TargetUserId, historySummary.Amount, historySummary.Type)
	if err != nil {
		return err
	}

	_, err = d.redis.Delete(ctx, fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, historySummary.UserId, historySummary.Type))
	if err != nil {
		return err
	}

	d.cache.Delete(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, historySummary.UserId, historySummary.Type))

	return nil
}

func (d domain) GrantBalanceByUserId(ctx context.Context, balance entity.Balance) (err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			d.db.Commit(tx)
		} else {
			d.db.Rollback(tx)
		}
	}()

	err = d.grantBalanceByUserId(ctx, tx, balance)
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, entity.History{
		Id:           uuid.NewString(),
		UserId:       balance.UserId,
		TargetUserId: balance.UserId,
		Amount:       balance.Amount,
		Type:         int(enum.CREDIT),
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
			d.db.Commit(tx)
		} else {
			d.db.Rollback(tx)
		}
	}()

	err = d.grantBalanceByUserId(ctx, tx, entity.Balance{
		UserId: req.ToUserId,
		Amount: req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.deductBalanceByUserId(ctx, tx, entity.Balance{
		UserId: req.UserId,
		Amount: req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, entity.History{
		Id:           uuid.NewString(),
		UserId:       req.ToUserId,
		TargetUserId: req.UserId,
		Amount:       req.Amount,
		Type:         int(enum.CREDIT),
		Notes:        fmt.Sprintf("Receive money from %s", req.UserId),
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, entity.History{
		Id:           uuid.NewString(),
		UserId:       req.UserId,
		TargetUserId: req.ToUserId,
		Amount:       req.Amount,
		Type:         int(enum.DEBIT),
		Notes:        fmt.Sprintf("Transfer money to %s", req.ToUserId),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d domain) GetLatestHistoryByUserId(ctx context.Context, userId string) (resp []entity.History, err error) {
	defer func() {
		if err == nil {
			for i := range resp {
				resp[i].NormalizeAmount()
			}
		}
	}()

	histories, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetLatestHistoryByUserId, userId), func() (interface{}, error) {
		var resp []entity.History
		histories, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, userId), time.Minute*5, func() (interface{}, error) {
			var respRedis []entity.History
			historiesStr, err := d.redis.Fetch(ctx, fmt.Sprintf(cacheKeyGetLatestHistoryByUserId, userId), time.Duration(time.Minute*30), func() (interface{}, error) {
				var history []entity.History
				err := d.db.SelectContextStmt(ctx, d.stmts.getLatestHistoryByUserId, &history, userId)
				if err != nil {
					return history, err
				}

				return history, nil
			})
			if err != nil {
				return respRedis, err
			}

			err = jsoniter.UnmarshalFromString(historiesStr, &respRedis)
			if err != nil {
				return respRedis, err
			}

			return respRedis, nil
		})
		if err != nil {
			return resp, err
		}

		return histories.Value().([]entity.History), nil
	})
	if err != nil {
		return resp, err
	}

	return histories.([]entity.History), nil
}

func (d domain) GetHistorySummaryByUserIdAndType(ctx context.Context, userId string, historyType int) (resp []entity.HistorySummary, err error) {
	historySummaries, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetHistorySummaryByUserIdAndType, userId, historyType), func() (interface{}, error) {
		var resp []entity.HistorySummary
		historySummaries, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, userId, historyType), time.Minute*5, func() (interface{}, error) {
			var respRedis []entity.HistorySummary
			historySummariesStr, err := d.redis.Fetch(ctx, fmt.Sprintf(cacheKeyGetHistorySummaryByUserIdAndType, userId, historyType), time.Duration(time.Minute*30), func() (interface{}, error) {
				var historySummary []entity.HistorySummary
				err := d.db.SelectContextStmt(ctx, d.stmts.getHistorySummaryByUserIdAndType, &historySummary, userId, historyType)
				if err != nil {
					return historySummary, err
				}

				return historySummary, nil
			})
			if err != nil {
				return respRedis, err
			}

			err = jsoniter.UnmarshalFromString(historySummariesStr, &respRedis)
			if err != nil {
				return respRedis, err
			}

			return respRedis, nil
		})
		if err != nil {
			return resp, err
		}

		return historySummaries.Value().([]entity.HistorySummary), nil
	})
	if err != nil {
		return resp, err
	}

	return historySummaries.([]entity.HistorySummary), nil
}
