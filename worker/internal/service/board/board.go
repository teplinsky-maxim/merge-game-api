package board

import (
	"context"
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	"merge-api/worker/pkg"
)

type CollectionBoardService struct {
	repo  *repo.CollectionBoard
	redis *repo.RedisCollectionBoard
}

func (r *CollectionBoardService) GetBoard(id uint) (board.Board[pkg.CollectionItem], error) {
	//TODO implement me
	panic("implement me")
}

func (r *CollectionBoardService) CreateBoard(w, h uint) (board.Board[pkg.CollectionItem], uint, error) {
	ctx := context.Background()
	createdBoard, boardId, err := (*r.repo).CreateBoard(ctx, w, h)
	if err != nil {
		return nil, 0, err
	}
	err = (*r.redis).CreateBoard(ctx, createdBoard, boardId)
	if err != nil {
		return nil, 0, err
	}
	return createdBoard, boardId, nil
}

func (r *CollectionBoardService) UpdateBoard(id uint, board *board.Board[pkg.CollectionItem]) error {
	//TODO implement me
	panic("implement me")
}

func (r *CollectionBoardService) DeleteBoard(id uint) error {
	//TODO implement me
	panic("implement me")
}

func NewCollectionBoardService(repo repo.CollectionBoard, redis repo.RedisCollectionBoard) *CollectionBoardService {
	return &CollectionBoardService{
		repo:  &repo,
		redis: &redis,
	}
}
