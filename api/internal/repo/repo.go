package repo

import (
	"context"
	"merge-api/api/internal/entity"
	"merge-api/api/internal/repo/repos/collection"
	task2 "merge-api/api/internal/repo/repos/task"
	"merge-api/shared/entity/task"
	"merge-api/shared/pkg/board"
	"merge-api/shared/pkg/database"
)

type Collection interface {
	GetCollections(ctx context.Context, offset, limit uint) ([]entity.Collection, error)
	GetCollection(ctx context.Context, collectionId uint) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, name string) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, items []entity.CollectionItem) ([]entity.CollectionItem, error)
}

type Task interface {
	CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task.Task, error)
}

type Repositories struct {
	Collection
	Task
}

func NewRepositories(database *database.Database) *Repositories {
	return &Repositories{
		Collection: collection.NewCollectionRepo(database),
		Task:       task2.NewTaskRepo(database),
	}
}
