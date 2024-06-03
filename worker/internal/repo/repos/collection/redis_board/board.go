package redis_board

import (
	"context"
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/service"
	"merge-api/worker/pkg/redis"
)

type Repo struct {
	redis *redis.Redis
}

func (r *Repo) GetBoard(ctx context.Context, id uint) (board.Board[service.CollectionItem], error) {
	panic("implement me")
}

func (r *Repo) CreateBoard(ctx context.Context, w, h uint) (board.Board[service.CollectionItem], uint, error) {
	panic("implement me")
}

func (r *Repo) UpdateBoard(ctx context.Context, id uint, board *board.Board[service.CollectionItem]) error {
	panic("implement me")
}

func (r *Repo) DeleteBoard(ctx context.Context, id uint) error {
	panic("implement me")
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
