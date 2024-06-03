package service

import (
	"context"
	"merge-api/api/internal/entity"
	"merge-api/api/internal/repo"
	"merge-api/api/internal/service/collection"
	task2 "merge-api/api/internal/service/task"
	"merge-api/api/pkg/board"
	taskEntity "merge-api/shared/entity/task"
	"merge-api/shared/pkg/rabbitmq"
)

type Collection interface {
	GetCollections(ctx context.Context, input *collection.GetCollectionsInput) ([]entity.Collection, error)
	GetCollection(ctx context.Context, input *collection.GetCollectionInput) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, input *collection.CreateCollectionInput) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, input *collection.CreateCollectionItemsInput) ([]entity.CollectionItem, error)
}

type Task interface {
	CreateTaskNewBoard(ctx context.Context, width, height board.SizeType) (taskEntity.Task, error)
	CreateTaskMoveItem(ctx context.Context /**/) (taskEntity.IDType, error)
	CreateTaskMergeItems(ctx context.Context /**/) (taskEntity.IDType, error)
	CreateTaskClickItem(ctx context.Context /**/) (taskEntity.IDType, error)
}

type Services struct {
	Collection Collection
	Task       Task
}
type Dependencies struct {
	Repositories repo.Repositories
	RabbitMQ     rabbitmq.RabbitMQ
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Collection: collection.NewCollectionService(deps.Repositories.Collection),
		Task:       task2.NewTaskService(deps.Repositories.Task, &deps.RabbitMQ),
	}
}
