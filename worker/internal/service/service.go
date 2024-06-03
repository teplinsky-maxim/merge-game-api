package service

import (
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	boardService "merge-api/worker/internal/service/board"
	"merge-api/worker/pkg/redis"
)

type CollectionItem interface{}

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
	Board[CollectionItem]
}

//type Task interface {
//	SetTaskStarted(taskId task.IDType) error
//	SetTaskDone(taskId task.IDType, result any) error
//	SetTaskFailed(taskId task.IDType) error
//}

type Services struct {
	Board CollectionBoard
	//Task  Task
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
