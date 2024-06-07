package domainauth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kevinsudut/wallet-system/app/entity"
)

func (d domain) InsertUser(ctx context.Context, user entity.User) (err error) {
	err = d.db.ExecContextStmt(ctx, d.stmts.insertUser, user.Id, user.Username)
	if err != nil {
		return err
	}

	json, err := jsoniter.MarshalToString(user)
	if err != nil {
		return err
	}

	_, err = d.redis.SetEx(ctx, fmt.Sprintf(cacheKeyGetUserById, user.Id), json, time.Minute*30)
	if err != nil {
		return err
	}

	_, err = d.redis.SetEx(ctx, fmt.Sprintf(cacheKeyGetUserByUsername, user.Username), json, time.Minute*30)
	if err != nil {
		return err
	}

	d.cache.Set(fmt.Sprintf(cacheKeyGetUserById, user.Id), user, time.Minute*5)
	d.cache.Set(fmt.Sprintf(cacheKeyGetUserByUsername, user.Username), user, time.Minute*5)

	return nil
}

func (d domain) GetUserById(ctx context.Context, id string) (resp entity.User, err error) {
	user, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetUserById, id), func() (interface{}, error) {
		var resp entity.User
		user, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetUserById, id), time.Minute*5, func() (interface{}, error) {
			var respRedis entity.User
			userStr, err := d.redis.Fetch(ctx, fmt.Sprintf(cacheKeyGetUserById, id), time.Duration(time.Minute*30), func() (interface{}, error) {
				var user entity.User
				err := d.db.GetContextStmt(ctx, d.stmts.getUserById, &user, id)
				if err != nil {
					return user, err
				}

				return user, nil
			})
			if err != nil {
				return respRedis, err
			}

			err = jsoniter.UnmarshalFromString(userStr, &respRedis)
			if err != nil {
				return respRedis, err
			}

			return respRedis, nil
		})
		if err != nil {
			return resp, err
		}

		return user.Value().(entity.User), nil
	})
	if err != nil {
		return resp, err
	}

	return user.(entity.User), nil
}

func (d domain) GetUserByUsername(ctx context.Context, username string) (resp entity.User, err error) {
	defer func() {
		if err == nil && resp.Id == "" {
			err = sql.ErrNoRows
		}
	}()

	user, err, _ := d.singleflight.DoSingleFlight(ctx, fmt.Sprintf(singleFlightKeyGetUserById, username), func() (interface{}, error) {
		var resp entity.User
		user, err := d.cache.Fetch(fmt.Sprintf(cacheKeyGetUserByUsername, username), time.Minute*5, func() (interface{}, error) {
			var respRedis entity.User
			userStr, err := d.redis.Fetch(ctx, fmt.Sprintf(cacheKeyGetUserByUsername, username), time.Duration(time.Minute*30), func() (interface{}, error) {
				var user entity.User
				err := d.db.GetContextStmt(ctx, d.stmts.getUserByUsername, &user, username)
				if err != nil && err != sql.ErrNoRows {
					return user, err
				}

				return user, nil
			})
			if err != nil {
				return respRedis, err
			}

			err = jsoniter.UnmarshalFromString(userStr, &respRedis)
			if err != nil {
				return respRedis, err
			}

			return respRedis, nil
		})
		if err != nil {
			return resp, err
		}

		return user.Value().(entity.User), nil
	})
	if err != nil {
		return resp, err
	}

	return user.(entity.User), nil
}
