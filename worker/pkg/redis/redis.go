package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"merge-api/shared/config"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(config *config.Config) (Redis, error) {
	dsn := fmt.Sprintf("%v:%v", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return Redis{}, err
	}
	return Redis{
		client: rdb,
	}, nil
}
