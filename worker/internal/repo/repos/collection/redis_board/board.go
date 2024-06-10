package redis_board

import (
	"context"
	"fmt"
	"merge-api/shared/pkg/board"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
)

const BoardRedisKey = "board"

type Repo struct {
	redis *redis.Redis
}

func makeRedisBoardKey(boardId uint) string {
	return fmt.Sprintf("%v:%v", BoardRedisKey, boardId)
}

func (r *Repo) CreateBoard(ctx context.Context, board board.Board[pkg.CollectionItem], boardId uint) error {
	redisKey := makeRedisBoardKey(boardId)
	err := r.redis.Client.HSet(ctx, redisKey, "", "").Err()
	return err
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
