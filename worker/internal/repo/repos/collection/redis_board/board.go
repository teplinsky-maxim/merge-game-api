package redis_board

import (
	"context"
	"merge-api/shared/pkg/board"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
)

type Repo struct {
	redis *redis.Redis
}

func (r *Repo) CreateBoard(ctx context.Context, board board.Board[pkg.CollectionItem], boardId uint) error {
	panic("implement me")
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
