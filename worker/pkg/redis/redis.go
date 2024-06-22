package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"merge-api/shared/config"
)

type RedisTest Redis
type Redis struct {
	Client *redis.Client
}

func NewRedis(config *config.Config) (Redis, error) {
	r, err := makeRedis(config, false)
	if err != nil {
		panic(err)
	}
	return r, nil
}

func NewTestRedis(config *config.Config) (RedisTest, error) {
	r, err := makeRedis(config, true)
	if err != nil {
		panic(err)
	}
	return RedisTest(r), nil
}

func makeRedis(config *config.Config, test bool) (Redis, error) {
	db := 0
	if test {
		db = 1
	}
	dsn := fmt.Sprintf("%v:%v", config.Redis.Host, config.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       db,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return Redis{}, err
	}
	return Redis{
		Client: rdb,
	}, nil
}

func (r *RedisTest) Clear() {
	err := r.Client.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}
