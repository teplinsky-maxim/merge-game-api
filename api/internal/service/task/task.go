package task

import (
	"context"
	"merge-api/api/internal/repo"
	"merge-api/api/pkg/board"
	taskEntity "merge-api/shared/entity/task"
	"merge-api/shared/pkg/rabbitmq"
	"merge-api/shared/pkg/rabbitmq/tasks"
)

type TaskService struct {
	repo repo.Task
	rmq  *rabbitmq.RabbitMQ
}

func (r *TaskService) CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (taskEntity.Task, error) {
	createdTask, err := r.repo.CreateTaskNewBoard(ctx, width, height)
	if err != nil {
		return taskEntity.Task{}, err
	}
	err = tasks.SendTask(r.rmq, createdTask)
	if err != nil {
		return taskEntity.Task{}, err
	}
	return createdTask, nil
}

func (r *TaskService) CreateTaskMoveItem(ctx context.Context) (taskEntity.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func (r *TaskService) CreateTaskMergeItems(ctx context.Context) (taskEntity.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func (r *TaskService) CreateTaskClickItem(ctx context.Context) (taskEntity.IDType, error) {
	//TODO implement me
	panic("implement me")
}

func NewTaskService(repo repo.Task, rmq *rabbitmq.RabbitMQ) *TaskService {
	return &TaskService{
		repo: repo,
		rmq:  rmq,
	}
}
