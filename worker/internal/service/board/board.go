package board

import (
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	"merge-api/worker/internal/service"
)

type CollectionBoardService struct {
	repo  *repo.CollectionBoard
	redis *repo.RedisCollectionBoard
}

func (r *CollectionBoardService) GetBoard(id uint) (board.Board[service.CollectionItem], error) {
	//TODO implement me
	panic("implement me")
}

func (r *CollectionBoardService) CreateBoard(w, h uint) (board.Board[service.CollectionItem], uint, error) {
	//TODO implement me
	panic("implement me")
}

func (r *CollectionBoardService) UpdateBoard(id uint, board *board.Board[service.CollectionItem]) error {
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
