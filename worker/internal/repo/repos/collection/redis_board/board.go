package redis_board

import (
	"context"
	"errors"
	"fmt"
	goRedis "github.com/redis/go-redis/v9"
	"merge-api/shared/pkg/board"
	"merge-api/worker/pkg"
	"merge-api/worker/pkg/redis"
	"strconv"
	"strings"
)

const BoardRedisKey = "board"

type PipelineTxType string

const TxKey = PipelineTxType("tx")

var BoardNotLoadedInRedisError = errors.New("board is not loaded in redis")
var BoardCellEmptyError = errors.New("board cell is empty")

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

	boardInRedis, err := r.redis.Client.Exists(ctx, boardKey).Result()
	if err != nil {
		return nil, err
	}
	if boardInRedis != 1 {
		return nil, BoardNotLoadedInRedisError
	}

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
	txWasCreatedHere := false
	pipeline, ok := ctx.Value(TxKey).(goRedis.Pipeliner)
	if !ok {
		txWasCreatedHere = true
		pipeline = r.redis.Client.TxPipeline()
	}

	coordinatesString := makeRedisCoordinatesString(w, h)
	boardKey := makeRedisBoardKey(id)

	boardInRedis, err := r.redis.Client.Exists(ctx, boardKey).Result()
	if err != nil {
		return err
	}
	if boardInRedis != 1 {
		return BoardNotLoadedInRedisError
	}

	_ = pipeline.HSet(ctx, boardKey, coordinatesString, collectionItem.ToRedisString())
	if txWasCreatedHere {
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) ClearCell(ctx context.Context, id, w, h uint) error {
	txWasCreatedHere := false
	pipeline, ok := ctx.Value(TxKey).(goRedis.Pipeliner)
	if !ok {
		txWasCreatedHere = true
		pipeline = r.redis.Client.TxPipeline()
	}

	coordinatesString := makeRedisCoordinatesString(w, h)
	boardKey := makeRedisBoardKey(id)

	boardInRedis, err := r.redis.Client.Exists(ctx, boardKey).Result()
	if err != nil {
		return err
	}
	if boardInRedis != 1 {
		return BoardNotLoadedInRedisError
	}

	_ = pipeline.HDel(ctx, boardKey, coordinatesString)
	if txWasCreatedHere {
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRedisBoardRepo(redis *redis.Redis) *Repo {
	return &Repo{
		redis: redis,
	}
}
