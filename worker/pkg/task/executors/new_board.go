package executors

import (
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/repo"
	"merge-api/worker/internal/service"
	"merge-api/worker/pkg/redis"
)

type NewBoardTaskExecutor struct {
	repo  *repo.Repositories
	redis *redis.Redis
}

func (n *NewBoardTaskExecutor) CanExecuteThisTask(t *task.Task) (bool, error) {
	return t.Type == task.NewBoard, nil
}

func (n *NewBoardTaskExecutor) Execute(t *task.Task) (any, error) {
	// Implement the task execution logic
	return nil, nil
}

func NewNewBoardTaskExecutor(dependencies service.Dependencies) *NewBoardTaskExecutor {
	return &NewBoardTaskExecutor{
		repo:  &dependencies.Repositories,
		redis: &dependencies.Redis,
	}
}
