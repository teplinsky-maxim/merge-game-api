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

type Repo struct {
	redis *redis.Redis
}

func makeRedisBoardKey(boardId uint) string {
	return fmt.Sprintf("%v:%v", BoardRedisKey, boardId)
}

func makeRedisCoordinatesString(w, h uint) string {
	return fmt.Sprintf("%v:%v", w, h)
}

func (r *Repo) CreateBoard(ctx context.Context, board board.Board[pkg.CollectionItem], boardId uint) error {
	redisKey := makeRedisBoardKey(boardId)
	err := r.redis.Client.HSet(ctx, redisKey, "", "").Err()
	return err
}

func (r *Repo) GetBoardByCoordinates(ctx context.Context, id, w, h uint) (pkg.CollectionItem, error) {
	coordinatesString := makeRedisCoordinatesString(w, h)
	boardKey := makeRedisBoardKey(id)

	result, err := r.redis.Client.HGet(ctx, boardKey, coordinatesString).Result()
	if err != nil {
		return nil, nil
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
	return pkg.NewCollectionItemImpl(uint(cIdInt), uint(cItemIdInt)), nil
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
