package domainbalance

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func (d domain) GetBalanceByUsername(ctx context.Context, username string) (resp Balance, err error) {
	balance, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetBalanceByUsername, username), time.Minute*5, func() (string, error) {
		var balance Balance
		err := d.stmts.getBalanceByUsername.GetContext(ctx, &balance, username)
		if err != nil {
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
}

func (d domain) grantBalanceByUsername(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	result, err := tx.ExecContext(ctx, queryGrantBalanceByUsername, balance.Username, balance.Amount)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetBalanceByUsername, balance.Username))

	return nil
}

func (d domain) deductBalanceByUsername(ctx context.Context, tx *sql.Tx, balance Balance) (err error) {
	result, err := tx.ExecContext(ctx, queryDeductBalanceByUsername, balance.Amount, balance.Username)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetBalanceByUsername, balance.Username))

	return nil
}

func (d domain) insertHistory(ctx context.Context, tx *sql.Tx, history History) (err error) {
	result, err := tx.ExecContext(ctx, queryInsertHistory, history.Id, history.Username, history.TargetUsername, history.Amount, history.Type, history.Notes)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	err = d.updateHistorySummary(ctx, tx, HistorySummary{
		Username:       history.Username,
		TargetUsername: history.TargetUsername,
		Amount:         history.Amount,
		Type:           history.Type,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d domain) updateHistorySummary(ctx context.Context, tx *sql.Tx, historySummary HistorySummary) (err error) {
	result, err := tx.ExecContext(ctx, queryUpdateHistorySummaryById, historySummary.GetId(), historySummary.Username, historySummary.TargetUsername, historySummary.Amount, historySummary.Type)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	return nil
}

func (d domain) GrantBalanceByUsername(ctx context.Context, balance Balance) (err error) {
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

	err = d.grantBalanceByUsername(ctx, tx, balance)
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:             uuid.NewString(),
		Username:       balance.Username,
		TargetUsername: balance.Username,
		Amount:         balance.Amount,
		Type:           int(CREDIT),
		Notes:          "Top-up money",
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

	err = d.grantBalanceByUsername(ctx, tx, Balance{
		Username: req.ToUsername,
		Amount:   req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.deductBalanceByUsername(ctx, tx, Balance{
		Username: req.Username,
		Amount:   req.Amount,
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:             uuid.NewString(),
		Username:       req.ToUsername,
		TargetUsername: req.Username,
		Amount:         req.Amount,
		Type:           int(CREDIT),
		Notes:          fmt.Sprintf("Receive money from %s", req.Username),
	})
	if err != nil {
		return err
	}

	err = d.insertHistory(ctx, tx, History{
		Id:             uuid.NewString(),
		Username:       req.Username,
		TargetUsername: req.ToUsername,
		Amount:         req.Amount,
		Type:           int(DEBIT),
		Notes:          fmt.Sprintf("Transfer money to %s", req.ToUsername),
	})
	if err != nil {
		return err
	}

	return nil
}
