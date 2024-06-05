package task

import (
	"context"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/database"
)

type Repo struct {
	database *database.Database
}

func (r *Repo) UpdateTask(ctx context.Context, taskId task.IDType, status task.Status, result any) error {
	panic("implement me")
}

func NewTaskRepo(database *database.Database) *Repo {
	return &Repo{
		database: database,
	}
}
