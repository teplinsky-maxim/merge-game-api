package repo

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/database"
	"merge-api/worker/internal/repo/repos/collection"
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
	GetBoardByCoordinates(ctx context.Context, id, w, h uint) (T, error)
	CreateBoard(ctx context.Context, w, h uint) (board.Board[T], uint, error)
	DeleteBoard(ctx context.Context, id uint) error

	UpdateCell(ctx context.Context, id, w, h uint, t T) error
	ClearCell(ctx context.Context, id, w, h uint) error
}

// CollectionBoard is board.go for collections
// you use it as high-level interface
type CollectionBoard interface {
	Board[pkg.CollectionItem]
}

// RedisBoard is an interface only for redis
type RedisBoard[T any] interface {
	CreateBoard(ctx context.Context, board board.Board[T], boardId uint) error
	GetBoardByCoordinates(ctx context.Context, id, w, h uint) (T, error)

	UpdateCell(ctx context.Context, id, w, h uint, t T) error
	ClearCell(ctx context.Context, id, w, h uint) error
}

// RedisCollectionBoard is to store board.go with collection in redis
type RedisCollectionBoard interface {
	RedisBoard[pkg.CollectionItem]
}

type Collection interface {
	GetNextCollectionItem(ctx context.Context, item pkg.CollectionItem) (pkg.CollectionItem, error)
	IsItemMergeable(ctx context.Context, item pkg.CollectionItem) (bool, error)
}

type Repositories struct {
	Task
	CollectionBoard
	RedisCollectionBoard
	Collection
}

func NewRepositories(database *database.Database, redis *redis.Redis) *Repositories {
	return &Repositories{
		Task:                 task2.NewTaskRepo(database),
		CollectionBoard:      board2.NewBoardRepo(database),
		RedisCollectionBoard: redis_board.NewRedisBoardRepo(redis),
		Collection:           collection.NewCollectionRepo(database),
	}
}
