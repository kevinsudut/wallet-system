package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type rdb struct {
	client *redis.Client
}

func Init() (RedisItf, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &rdb{
		client: conn,
	}, nil
}
