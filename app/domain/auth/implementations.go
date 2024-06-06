package domainauth

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (d domain) InsertUser(ctx context.Context, user User) (err error) {
	err = d.db.ExecContextStmt(ctx, d.stmts.insertUser, user.Id, user.Username)
	if err != nil {
		return err
	}

	d.cache.Set(fmt.Sprintf(memcacheKeyGetUserById, user.Id), user, time.Minute*5)
	d.cache.Set(fmt.Sprintf(memcacheKeyGetUserByUsername, user.Username), user, time.Minute*5)

	return nil
}

func (d domain) GetUserById(ctx context.Context, id string) (resp User, err error) {
	user, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetUserById, id), func() (interface{}, error) {
		var resp User
		user, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetUserById, id), time.Minute*5, func() (interface{}, error) {
			var user User
			err := d.db.GetContextStmt(ctx, d.stmts.getUserById, &user, id)
			if err != nil {
				return user, err
			}

			return user, nil
		})
		if err != nil {
			return resp, err
		}

		return user.Value().(User), nil
	})
	if err != nil {
		return resp, err
	}

	return user.(User), nil
}

func (d domain) GetUserByUsername(ctx context.Context, username string) (resp User, err error) {
	defer func() {
		if err == nil && resp.Id == "" {
			err = sql.ErrNoRows
		}
	}()

	user, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetUserById, username), func() (interface{}, error) {
		var resp User
		user, err := d.cache.Fetch(fmt.Sprintf(memcacheKeyGetUserByUsername, username), time.Minute*5, func() (interface{}, error) {
			var user User
			err := d.db.GetContextStmt(ctx, d.stmts.getUserByUsername, &user, username)
			if err != nil && err != sql.ErrNoRows {
				return user, err
			}

			return user, nil
		})
		if err != nil {
			return resp, err
		}

		return user.Value().(User), nil
	})
	if err != nil {
		return resp, err
	}

	return user.(User), nil
}
