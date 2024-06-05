package service

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	boardService "merge-api/worker/internal/service/board"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
)

// Board is high-level interface
type Board[T any] interface {
	GetBoard(id uint) (board.Board[T], error)            // Get board.go, you can get it from cache or database or whatever you want
	CreateBoard(w, h uint) (board.Board[T], uint, error) // Create board.go
	UpdateBoard(id uint, board *board.Board[T]) error    // Update board.go
	DeleteBoard(id uint) error                           // Delete board.go
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

type Services struct {
	Board CollectionBoard
	Task  Task
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
	}
}
