package board

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/board/inmemory"
	"merge-api/shared/pkg/database"
	"merge-api/worker/pkg"
)

type Repo struct {
	database *database.Database
}

func (r *Repo) GetBoard(ctx context.Context, id uint) (board.Board[pkg.CollectionItem], error) {
	panic("implement me")
}

func (r *Repo) CreateBoard(ctx context.Context, w, h uint) (board.Board[pkg.CollectionItem], uint, error) {
	tx, err := r.database.DB.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}

	stmt := sq.
		Insert("boards").
		Columns("width", "height").
		Values(w, h).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	result := inmemory.NewBoard[pkg.CollectionItem](w, h)
	var id uint
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, 0, err
		} else {
			err = rows.Scan(&id)
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, 0, err
	}
	return &result, id, nil
}

func (r *Repo) UpdateBoard(ctx context.Context, id uint, board *board.Board[pkg.CollectionItem]) error {
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
