package task

import (
	"context"
	task2 "merge-api/internal/entity/task"
	"merge-api/internal/repo"
	"merge-api/pkg/board"
	"merge-api/pkg/rabbitmq"
	"merge-api/pkg/rabbitmq/tasks"
	"merge-api/pkg/task"
)

type TaskService struct {
	repo repo.Task
	rmq  *rabbitmq.RabbitMQ
}

func (r *TaskService) CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task2.Task, error) {
	createdTask, err := r.repo.CreateTaskNewBoard(ctx, width, height)
	if err != nil {
		return task2.Task{}, err
	}
	err = tasks.SendTask(r.rmq, createdTask)
	if err != nil {
		return task2.Task{}, err
	}
	return createdTask, nil
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

func NewTaskService(repo repo.Task, rmq *rabbitmq.RabbitMQ) *TaskService {
	return &TaskService{
		repo: repo,
		rmq:  rmq,
	}
}
