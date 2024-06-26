package task

import (
	"context"
	"merge-api/api/internal/repo"
	taskEntity "merge-api/shared/entity/task"
	"merge-api/shared/pkg/rabbitmq"
	"merge-api/shared/pkg/rabbitmq/tasks"
)

type TaskService struct {
	repo repo.Task
	rmq  *rabbitmq.RabbitMQ
}

func (r *TaskService) CreateTaskNewBoard(ctx context.Context, width, height uint) (taskEntity.Task, error) {
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

func (r *TaskService) CreateTaskMoveItem(ctx context.Context, boardId, w1, h1, w2, h2 uint) (taskEntity.Task, error) {
	createdTask, err := r.repo.CreateTaskMoveItem(ctx, boardId, w1, h1, w2, h2)
	if err != nil {
		return taskEntity.Task{}, err
	}
	err = tasks.SendTask(r.rmq, createdTask)
	if err != nil {
		return taskEntity.Task{}, err
	}
	return createdTask, nil
}

func (r *TaskService) CreateTaskMergeItems(ctx context.Context, boardId, w1, h1, w2, h2 uint) (taskEntity.Task, error) {
	createdTask, err := r.repo.CreateTaskMergeItems(ctx, boardId, w1, h1, w2, h2)
	if err != nil {
		return taskEntity.Task{}, err
	}
	err = tasks.SendTask(r.rmq, createdTask)
	if err != nil {
		return taskEntity.Task{}, err
	}
	return createdTask, nil
}

func (r *TaskService) CreateTaskClickItem(ctx context.Context, boardId, w1, h1 uint) (taskEntity.Task, error) {
	createdTask, err := r.repo.CreateTaskClickItem(ctx, boardId, w1, h1)
	if err != nil {
		return taskEntity.Task{}, err
	}
	err = tasks.SendTask(r.rmq, createdTask)
	if err != nil {
		return taskEntity.Task{}, err
	}
	return createdTask, nil
}

func NewTaskService(repo repo.Task, rmq *rabbitmq.RabbitMQ) *TaskService {
	return &TaskService{
		repo: repo,
		rmq:  rmq,
	}
}
