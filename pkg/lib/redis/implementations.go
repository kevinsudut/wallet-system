package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kevinsudut/wallet-system/pkg/lib/log"
	"github.com/redis/go-redis/v9"
)

func (r rdb) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r rdb) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.client.SetEx(ctx, key, value, expiration).Result()
}

func (r rdb) Fetch(ctx context.Context, key string, expiration time.Duration, fetch func() (interface{}, error)) (string, error) {
	resp, err := r.Get(ctx, key)
	if err == nil {
		fmt.Println("GET FROM REDIS", key)
		return resp, nil
	} else if !errors.Is(err, redis.Nil) {
		log.Errorln("Redis.Fetch.Get", key, err)
	}

	fetchResp, err := fetch()
	if err != nil {
		return "", err
	}

	json, err := jsoniter.MarshalToString(fetchResp)
	if err != nil {
		return "", err
	}

	_, err = r.SetEx(ctx, key, json, expiration)
	if err != nil {
		log.Errorln("Redis.Fetch.SetEx", key, err)
	}

	fmt.Println("SET TO REDIS", key)

	return json, nil
}
