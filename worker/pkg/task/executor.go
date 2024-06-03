package task

import (
	"errors"
	"merge-api/shared/entity/task"
)

var (
	UnknownTask = errors.New("unknown task")
)

type Executor interface {
	CanExecuteThisTask(t *task.Task) (bool, error)
	Execute(t *task.Task) (any, error)
}

type ExecutorsManager struct {
	executors []Executor
}

func NewTaskExecutorsManager(executors []Executor) *ExecutorsManager {
	return &ExecutorsManager{executors: executors}
}

func (manager *ExecutorsManager) FindExecutor(t *task.Task) (Executor, error) {
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

func (manager *ExecutorsManager) ExecuteTask(t *task.Task) (any, error) {
	executor, err := manager.FindExecutor(t)
	if err != nil {
		return nil, err
	}
	return executor.Execute(t)
}
