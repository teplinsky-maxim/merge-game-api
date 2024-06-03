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
	GetBoard(id uint) (board.Board[T], error)            // Get board, you can get it from cache or database or whatever you want
	CreateBoard(w, h uint) (board.Board[T], uint, error) // Create board
	UpdateBoard(id uint, board *board.Board[T]) error    // Update board
	DeleteBoard(id uint) error                           // Delete board
}

// CollectionBoard is board for collections
// you use it as high-level interface
type CollectionBoard interface {
	Board[CollectionItem]
}

type Services struct {
	Board CollectionBoard
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
