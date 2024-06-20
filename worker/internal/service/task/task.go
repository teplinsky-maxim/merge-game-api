package task

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/worker/internal/repo"
)

type TaskService struct {
	repo *repo.Task
}

func (r *TaskService) SetTaskStarted(ctx context.Context, taskId task.IDType) error {
	err := (*r.repo).UpdateTask(ctx, taskId, task.Running, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskService) SetTaskDone(ctx context.Context, taskId task.IDType, result any) error {
	err := (*r.repo).UpdateTask(ctx, taskId, task.Done, result)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskService) SetTaskFailed(ctx context.Context, taskId task.IDType) error {
	err := (*r.repo).UpdateTask(ctx, taskId, task.Failed, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewTaskService(repo repo.Task) *TaskService {
	return &TaskService{
		repo: &repo,
	}
}
