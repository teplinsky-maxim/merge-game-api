package repo

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo/repos/collection"
	task2 "merge-api/internal/repo/repos/task"
	"merge-api/pkg/board"
	"merge-api/pkg/database"
	"merge-api/pkg/task"
)

type Collection interface {
	GetCollections(ctx context.Context, offset, limit uint) ([]entity.Collection, error)
	GetCollection(ctx context.Context, collectionId uint) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, name string) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, items []entity.CollectionItem) ([]entity.CollectionItem, error)
}

type Task interface {
	CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task.IDType, error)
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
