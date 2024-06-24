package executors

import (
	"encoding/json"
	"errors"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/service"
)

var (
	IncorrectBoardDimensions = errors.New("incorrect board dimension")
	BoardIsTooBig            = errors.New("board is too big")
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
	if args.Width == 0 || args.Height == 0 {
		return nil, IncorrectBoardDimensions
	}
	if args.Width*args.Height > 1000 {
		return nil, BoardIsTooBig
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
