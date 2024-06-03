package task

import (
	"errors"
	"merge-api/shared/entity/task"
)

var (
	UnknownTask = errors.New("unknown task")
)

type TaskExecutor interface {
	CanExecuteThisTask(t *task.Task) (bool, error)
	Execute(t *task.Task) (any, error)
}

type TaskExecutorsManager struct {
	executors []TaskExecutor
}

func NewTaskExecutorsManager(executors []TaskExecutor) *TaskExecutorsManager {
	return &TaskExecutorsManager{executors: executors}
}

func (manager *TaskExecutorsManager) FindExecutor(t *task.Task) (TaskExecutor, error) {
	for _, executor := range manager.executors {
		canExecute, err := executor.CanExecuteThisTask(t)
		if err != nil {
			return nil, err
		}
		if canExecute {
			return executor, nil
		}
	}
	return nil, UnknownTask
}

func (manager *TaskExecutorsManager) ExecuteTask(t *task.Task) (any, error) {
	executor, err := manager.FindExecutor(t)
	if err != nil {
		return nil, err
	}
	return executor.Execute(t)
}

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
