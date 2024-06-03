package service

import (
	"context"
	"merge-api/shared/entity/task"
)

type Task interface {
	Execute(ctx context.Context, task *task.Task) error
}
