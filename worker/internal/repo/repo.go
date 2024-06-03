package repo

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
)

type Task interface {
	UpdateTask(ctx context.Context, status task.Status, result any) error
}

type Repositories struct {
	Task
}

func NewRepositories(database *database.Database) *Repositories {
	return &Repositories{}
}
