package domainbalance

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func (d domain) GetBalanceByUsername(ctx context.Context, username string) (resp Balance, err error) {
	balance, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetBalanceByUsername, username), time.Minute*5, func() (string, error) {
		var balance Balance
		err := d.db.GetContext(ctx, &balance, queryGetBalanceByUsername, username)
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

func (d domain) GrantBalanceByUsername(ctx context.Context, balance Balance) (err error) {
	result, err := d.db.ExecContext(ctx, queryGrantBalanceByUsername, balance.Username, balance.Amount)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed update to database")
	}

	d.cache.Delete(fmt.Sprintf(memcacheKeyGetBalanceByUsername, balance.Username))

	return nil
}
