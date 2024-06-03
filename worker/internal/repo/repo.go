package repo

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/database"
	"merge-api/worker/internal/service"
)

type Task interface {
	UpdateTask(ctx context.Context, status task.Status, result any) error
}

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
	Board[service.CollectionItem]
}

// RedisBoard is an interface only for redis
type RedisBoard[T any] interface {
	Board[T]
}

// RedisCollectionBoard is to store board with collection in redis
type RedisCollectionBoard interface {
	RedisBoard[service.CollectionItem]
}

type Repositories struct {
	Task
	CollectionBoard
	RedisCollectionBoard
}

func NewRepositories(database *database.Database) *Repositories {
	return &Repositories{}
}
