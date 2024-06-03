package board

import (
	"context"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/database"
	"merge-api/worker/internal/service"
)

type Repo struct {
	database *database.Database
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

func NewBoardRepo(database *database.Database) *Repo {
	return &Repo{
		database: database,
	}
}
