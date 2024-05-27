package service

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo"
	"merge-api/internal/service/collection"
	task2 "merge-api/internal/service/task"
	"merge-api/pkg/board"
	"merge-api/pkg/task"
)

type Collection interface {
	GetCollections(ctx context.Context, input *collection.GetCollectionsInput) ([]entity.Collection, error)
	GetCollection(ctx context.Context, input *collection.GetCollectionInput) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, input *collection.CreateCollectionInput) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, input *collection.CreateCollectionItemsInput) ([]entity.CollectionItem, error)
}

type Task interface {
	CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (task.IDType, error)
	CreateTaskMoveItem(ctx context.Context /**/) (task.IDType, error)
	CreateTaskMergeItems(ctx context.Context /**/) (task.IDType, error)
	CreateTaskClickItem(ctx context.Context /**/) (task.IDType, error)
}

type Services struct {
	Collection Collection
	Task       Task
}
type Dependencies struct {
	Repositories repo.Repositories
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Collection: collection.NewCollectionService(deps.Repositories.Collection),
		Task:       task2.NewTaskService(deps.Repositories.Task),
	}
}
