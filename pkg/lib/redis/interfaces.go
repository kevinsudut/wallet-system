package redis

import (
	"context"
	"time"
)

type RedisItf interface {
	Get(ctx context.Context, key string) (string, error)
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Delete(ctx context.Context, keys ...string) (int64, error)
	Fetch(ctx context.Context, key string, expiration time.Duration, fetch func() (interface{}, error)) (string, error)
}
