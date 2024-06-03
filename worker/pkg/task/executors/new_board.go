package executors

import (
	"encoding/json"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
)

type NewBoardTaskExecutor struct {
	service *service.CollectionBoard
}

func (n *NewBoardTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.NewBoard, nil
}

func (n *NewBoardTaskExecutor) Execute(t *task.Task) (any, error) {
	args := new(task.NewBoardTaskArgs)
	err := json.Unmarshal(t.Args, args)
	if err != nil {
		return nil, err
	}
	_, boardId, err := (*n.service).CreateBoard(args.Width, args.Height)
	if err != nil {
		return nil, err
	}
	result := new(task.NewBoardTaskResult)
	result.BoardID = boardId
	return result, nil
}

func NewNewBoardTaskExecutor(boardService service.CollectionBoard) *NewBoardTaskExecutor {
	return &NewBoardTaskExecutor{
		service: &boardService,
	}
}
