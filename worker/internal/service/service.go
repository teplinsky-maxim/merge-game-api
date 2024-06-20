package service

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	boardService "merge-api/worker/internal/service/board"
	"merge-api/worker/internal/service/collection"
	taskService "merge-api/worker/internal/service/task"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
)

// Board is high-level interface
type Board[T any] interface {
	GetBoard(id uint) (board.Board[T], error)            // Get board.go, you can get it from cache or database or whatever you want
	GetBoardByCoordinates(id uint, w, h uint) (T, error) // Get board.go, you can get it from cache or database or whatever you want
	CreateBoard(w, h uint) (board.Board[T], uint, error) // Create board.go
	DeleteBoard(id uint) error                           // Delete board.go

	UpdateCell(id, w, h uint, t T) error
	ClearCell(id, w, h uint) error
}

// CollectionBoard is board.go for collections
// you use it as high-level interface
type CollectionBoard interface {
	Board[pkg.CollectionItem]
}

type Task interface {
	SetTaskStarted(ctx context.Context, taskId task.IDType) error
	SetTaskDone(ctx context.Context, taskId task.IDType, result any) error
	SetTaskFailed(ctx context.Context, taskId task.IDType) error
}

type Collection interface {
	GetNextCollectionItem(ctx context.Context, item pkg.CollectionItem) (pkg.CollectionItem, error)
	IsItemMergeable(ctx context.Context, item pkg.CollectionItem) (bool, error)
}

type Services struct {
	Board      CollectionBoard
	Task       Task
	Collection Collection
}

type Dependencies struct {
	Repositories repo.Repositories
	Redis        redis.Redis
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Board: boardService.NewCollectionBoardService(
			deps.Repositories.CollectionBoard,
			deps.Repositories.RedisCollectionBoard,
		),
		Task: taskService.NewTaskService(
			deps.Repositories.Task,
		),
		Collection: collection.NewCollectionService(
			deps.Repositories.Collection,
		),
	}
}
