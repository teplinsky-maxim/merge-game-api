package redis_board

import (
	"context"
	"errors"
	"fmt"
	"merge-api/shared/pkg/board"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
	"strconv"
	"strings"
)

const BoardRedisKey = "board"
const BoardRedisFieldNameForSize = "size"

var BoardNotLoadedInRedisError = errors.New("board is not loaded in redis")
var BoardCellEmptyError = errors.New("board cell is empty")
var CoordinatesOutOfBoundError = errors.New("coordinates out of bound error")

type Repo struct {
	redis *redis.Redis
}

func makeRedisBoardKey(boardId uint) string {
	return fmt.Sprintf("%v:%v", BoardRedisKey, boardId)
}

func makeRedisBoardFieldNameForSize() string {
	return BoardRedisFieldNameForSize
}

func makeRedisCoordinatesString(w, h uint) string {
	return fmt.Sprintf("%v:%v", w, h)
}

func makeRedisBoardSize(w, h uint) string {
	return fmt.Sprintf("%v:%v", w, h)
}

func (r *Repo) keyInRedis(ctx context.Context, key string) (bool, error) {
	keyExists, err := r.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return keyExists == 1, nil
}

func (r *Repo) boardInRedis(ctx context.Context, boardId uint) (bool, error) {
	boardKey := makeRedisBoardKey(boardId)
	return r.keyInRedis(ctx, boardKey)
}

func (r *Repo) getMaxBoardSize(ctx context.Context, boardId uint) (uint, uint, error) {
	boardKey := makeRedisBoardKey(boardId)
	boardInRedis, err := r.boardInRedis(ctx, boardId)
	if err != nil {
		return 0, 0, err
	}
	if !boardInRedis {
		return 0, 0, BoardNotLoadedInRedisError
	}

	sizeResult := r.redis.Client.HGet(ctx, boardKey, makeRedisBoardFieldNameForSize())
	boardSize, err := sizeResult.Result()
	boardSizeWH := strings.Split(boardSize, ":")
	if len(boardSizeWH) != 2 {
		return 0, 0, errors.New("somehow error key in redis is wrong")
	}
	boardW, boardH := boardSizeWH[0], boardSizeWH[1]

	boardWNumber, err := strconv.ParseUint(boardW, 10, 32)
	if err != nil {
		return 0, 0, err
	}
	boardHNumber, err := strconv.ParseUint(boardH, 10, 32)
	if err != nil {
		return 0, 0, err
	}
	return uint(boardWNumber), uint(boardHNumber), nil
}

func (r *Repo) isOutOfBounds(w, h, maxW, maxH uint) bool {
	return w > maxW || h > maxH
}

func (r *Repo) CreateBoard(ctx context.Context, board board.Board[pkg.CollectionItem], boardId uint) error {
	boardKey := makeRedisBoardKey(boardId)
	w := board.Width()
	h := board.Height()

	err := r.redis.Client.HSet(ctx, boardKey, makeRedisBoardFieldNameForSize(), makeRedisBoardSize(w, h)).Err()
	return err
}

func (r *Repo) GetBoardByCoordinates(ctx context.Context, id, w, h uint) (pkg.CollectionItem, error) {
	maxW, maxH, err := r.getMaxBoardSize(ctx, id)
	if err != nil {
		if !errors.Is(err, BoardNotLoadedInRedisError) {
			return nil, err
		}
	}
	if r.isOutOfBounds(w, h, maxW, maxH) {
		return nil, CoordinatesOutOfBoundError
	}

	boardKey := makeRedisBoardKey(id)
	coordinatesString := makeRedisCoordinatesString(w, h)

	result, err := r.redis.Client.HGet(ctx, boardKey, coordinatesString).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, BoardCellEmptyError
		}
		return nil, err
	}

	splitResult := strings.Split(result, ":")
	if len(splitResult) != 2 {
		return nil, errors.New("somehow error key in redis is wrong")
	}
	collectionId, collectionItemId := splitResult[0], splitResult[1]
	cIdInt, err := strconv.Atoi(collectionId)
	if err != nil {
		return nil, err
	}
	cItemIdInt, err := strconv.Atoi(collectionItemId)
	if err != nil {
		return nil, err
	}
	collectionItem := pkg.NewCollectionItemImpl(uint(cIdInt), uint(cItemIdInt))
	return &collectionItem, nil
}

func (r *Repo) UpdateCell(ctx context.Context, id, w, h uint, collectionItem pkg.CollectionItem) error {
	maxW, maxH, err := r.getMaxBoardSize(ctx, id)
	if err != nil {
		if !errors.Is(err, BoardNotLoadedInRedisError) {
			return err
		}
	}
	if r.isOutOfBounds(w, h, maxW, maxH) {
		return CoordinatesOutOfBoundError
	}

	boardKey := makeRedisBoardKey(id)
	coordinatesString := makeRedisCoordinatesString(w, h)

	err = r.redis.Client.HSet(ctx, boardKey, coordinatesString, collectionItem.ToRedisString()).Err()
	return err
}

func (r *Repo) ClearCell(ctx context.Context, id, w, h uint) error {
	maxW, maxH, err := r.getMaxBoardSize(ctx, id)
	if err != nil {
		if !errors.Is(err, BoardNotLoadedInRedisError) {
			return err
		}
	}
	if r.isOutOfBounds(w, h, maxW, maxH) {
		return CoordinatesOutOfBoundError
	}

	boardKey := makeRedisBoardKey(id)
	coordinatesString := makeRedisCoordinatesString(w, h)

	err = r.redis.Client.HDel(ctx, boardKey, coordinatesString).Err()
	return err
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
