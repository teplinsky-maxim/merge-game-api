package executors

import (
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
	// Implement the task execution logic
	return nil, nil
}

func NewNewBoardTaskExecutor(boardService service.CollectionBoard) *NewBoardTaskExecutor {
	return &NewBoardTaskExecutor{
		service: &boardService,
	}
}
