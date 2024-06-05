package domainauth

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func (d domain) InsertUser(ctx context.Context, user User) (err error) {
	result, err := d.stmts.insertUser.ExecContext(ctx, user.Id, user.Username)
	if err != nil {
		return err
	}

	if row, err := result.RowsAffected(); err != nil || row <= 0 {
		return fmt.Errorf("failed insert to database")
	}

	userStr, err := jsoniter.MarshalToString(user)
	if err == nil {
		d.cache.Set(fmt.Sprintf(memcacheKeyGetUserByUsername, user.Username), userStr, time.Minute*5)
	}

	return nil
}

func (d domain) GetUserByUsername(ctx context.Context, username string) (resp User, err error) {
	user, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetUserByUsername, username), time.Minute*5, func() (string, error) {
		var user User
		err := d.stmts.getUserByUsername.GetContext(ctx, &user, username)
		if err != nil {
			return "", err
		}

		return jsoniter.MarshalToString(user)
	})
	if err != nil {
		return resp, err
	}

	err = jsoniter.UnmarshalFromString(user.Value(), &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
