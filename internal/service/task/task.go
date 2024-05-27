package task

import (
	"context"
	"merge-api/internal/repo"
	"merge-api/pkg/board"
	"merge-api/pkg/task"
)

type TaskService struct {
	repo repo.Task
}

func (r *TaskService) CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task.IDType, error) {
	taskId, err := r.repo.CreateTaskNewBoard(ctx, width, height)
	if err != nil {
		return [16]byte{}, err
	}
	return taskId, nil
}

func (r *TaskService) CreateTaskMoveItem(ctx context.Context) (task.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func (r *TaskService) CreateTaskMergeItems(ctx context.Context) (task.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func (r *TaskService) CreateTaskClickItem(ctx context.Context) (task.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func NewTaskService(repo repo.Task) *TaskService {
	return &TaskService{repo: repo}
}
