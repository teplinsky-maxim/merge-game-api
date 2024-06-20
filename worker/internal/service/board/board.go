package board

import (
	"context"
	"errors"
	"merge-api/shared/pkg/board"
	"merge-api/worker/internal/repo"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
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

func (r *CollectionBoardService) GetBoardByCoordinates(id uint, w, h uint) (pkg.CollectionItem, error) {
	ctx := context.Background()
	atCoords, err := (*r.redis).GetBoardByCoordinates(ctx, id, w, h)
	if err != nil {
		if errors.Is(err, redis_board.BoardNotLoadedInRedisError) {
			atCoords, err = (*r.repo).GetBoardByCoordinates(ctx, id, w, h)
			return atCoords, err
		} else if errors.Is(err, redis_board.BoardCellEmptyError) {
			return nil, err
		} else {
			return nil, err
		}
	}
	return atCoords, nil
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

func (r *CollectionBoardService) UpdateCell(id, w, h uint, collectionItem pkg.CollectionItem) error {
	ctx := context.Background()
	// TODO: обеспечить синхронизацию состояний в redis и postgresql путем контекстов и транзакций
	err := (*r.redis).UpdateCell(ctx, id, w, h, collectionItem)
	if err != nil {
		if !errors.Is(err, redis_board.BoardNotLoadedInRedisError) {
			return err
		}
	}
	return (*r.repo).UpdateCell(ctx, id, w, h, collectionItem)
}

func (r *CollectionBoardService) ClearCell(id, w, h uint) error {
	ctx := context.Background()
	// TODO: обеспечить синхронизацию состояний в redis и postgresql путем контекстов и транзакций
	err := (*r.redis).ClearCell(ctx, id, w, h)
	if err != nil {
		if !errors.Is(err, redis_board.BoardNotLoadedInRedisError) {
			return err
		}
	}
	return (*r.repo).ClearCell(ctx, id, w, h)
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
