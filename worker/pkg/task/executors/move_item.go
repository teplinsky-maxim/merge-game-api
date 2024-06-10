package executors

import (
	"encoding/json"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
)

type MoveItemTaskExecutor struct {
	service *service.CollectionBoard
}

func (n *MoveItemTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.NewBoard, nil
}

func (n *MoveItemTaskExecutor) Execute(t *task.Task) (any, error) {
	args := new(task.MoveItemTaskArgs)
	err := json.Unmarshal(t.Args, args)
	if err != nil {
		return nil, err
	}
	atCoords2, err := (*n.service).GetBoardByCoordinates(args.BoardID, args.W1, args.H1)
	if err != nil {
		return nil, err
	}
	atCoords1, err := (*n.service).GetBoardByCoordinates(args.BoardID, args.W2, args.H2)
	if err != nil {
		return nil, err
	}
	// TODO: if both have the same item -> do nothing
	// else swap
	println(atCoords1, atCoords2)
	return nil, nil
}

func NewMoveItemTaskExecutor(boardService service.CollectionBoard) *MoveItemTaskExecutor {
	return &MoveItemTaskExecutor{
		service: &boardService,
	}
}
