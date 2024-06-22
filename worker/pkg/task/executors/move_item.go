package executors

import (
	"encoding/json"
	"errors"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/repo/repos/collection/redis_board"
	"merge-api/worker/internal/service"
)

var CoordinatesAreTheSameError = errors.New("coordinates are the same")

type MoveItemTaskExecutor struct {
	service *service.CollectionBoard
}

func (n *MoveItemTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.MoveItem, nil
}

func (n *MoveItemTaskExecutor) Execute(t *task.Task) (any, error) {
	args := new(task.MoveItemTaskArgs)
	err := json.Unmarshal(t.Args, args)
	if err != nil {
		return nil, err
	}

	if args.W1 == args.W2 && args.H1 == args.H2 {
		return nil, CoordinatesAreTheSameError
	}

	atCoords1, err1 := (*n.service).GetBoardByCoordinates(args.BoardID, args.W1, args.H1)
	if err1 != nil {
		if !errors.Is(err1, redis_board.BoardCellEmptyError) {
			return nil, err1
		}
	}
	atCoords2, err2 := (*n.service).GetBoardByCoordinates(args.BoardID, args.W2, args.H2)
	if err2 != nil {
		if !errors.Is(err2, redis_board.BoardCellEmptyError) {
			return nil, err2
		}
	}
	if errors.Is(err1, redis_board.BoardCellEmptyError) && !errors.Is(err2, redis_board.BoardCellEmptyError) {
		err = (*n.service).UpdateCell(args.BoardID, args.W1, args.H1, atCoords2)
		if err != nil {
			return nil, err
		}
		err = (*n.service).ClearCell(args.BoardID, args.W2, args.H2)
		if err != nil {
			return nil, err
		}
	} else if !errors.Is(err1, redis_board.BoardCellEmptyError) && errors.Is(err2, redis_board.BoardCellEmptyError) {
		err = (*n.service).UpdateCell(args.BoardID, args.W2, args.H2, atCoords1)
		if err != nil {
			return nil, err
		}
		err = (*n.service).ClearCell(args.BoardID, args.W1, args.H1)
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err1, redis_board.BoardCellEmptyError) && errors.Is(err2, redis_board.BoardCellEmptyError) {
		return nil, nil
	} else if !atCoords1.Equal(atCoords2) {
		err = (*n.service).UpdateCell(args.BoardID, args.W2, args.H2, atCoords1)
		if err != nil {
			return nil, err
		}
		err = (*n.service).UpdateCell(args.BoardID, args.W1, args.H1, atCoords2)
		if err != nil {
			return nil, err
		}
	}
	//TODO: return a result
	return nil, nil
}

func NewMoveItemTaskExecutor(boardService service.CollectionBoard) *MoveItemTaskExecutor {
	return &MoveItemTaskExecutor{
		service: &boardService,
	}
}
