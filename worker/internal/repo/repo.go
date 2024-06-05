package repo

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/database"
	board2 "merge-api/worker/internal/repo/repos/collection/board"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
	task2 "merge-api/worker/internal/repo/repos/task"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
)

type Task interface {
	UpdateTask(ctx context.Context, taskId task.IDType, status task.Status, result any) error
}

// Board is high-level interface
type Board[T any] interface {
	GetBoard(ctx context.Context, id uint) (board.Board[T], error)
	CreateBoard(ctx context.Context, w, h uint) (board.Board[T], uint, error)
	UpdateBoard(ctx context.Context, id uint, board *board.Board[T]) error
	DeleteBoard(ctx context.Context, id uint) error
}

// CollectionBoard is board.go for collections
// you use it as high-level interface
type CollectionBoard interface {
	Board[pkg.CollectionItem]
}

// RedisBoard is an interface only for redis
type RedisBoard[T any] interface {
	CreateBoard(ctx context.Context, board board.Board[T], boardId uint) error
}

// RedisCollectionBoard is to store board.go with collection in redis
type RedisCollectionBoard interface {
	RedisBoard[pkg.CollectionItem]
}

type Repositories struct {
	Task
	CollectionBoard
	RedisCollectionBoard
}

func NewRepositories(database *database.Database, redis *redis.Redis) *Repositories {
	return &Repositories{
		Task:                 task2.NewTaskRepo(database),
		CollectionBoard:      board2.NewBoardRepo(database),
		RedisCollectionBoard: redis_board.NewRedisBoardRepo(redis),
	}
}
