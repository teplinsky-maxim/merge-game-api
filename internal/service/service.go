package service

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo"
	"merge-api/internal/service/collection"
)

type Collection interface {
	GetCollections(ctx context.Context, input *collection.GetCollectionsInput) ([]entity.Collection, error)
	GetCollection(ctx context.Context, input *collection.GetCollectionInput) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, input *collection.CreateCollectionInput) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, input *collection.CreateCollectionItemsInput) ([]entity.CollectionItem, error)
}

type Services struct {
	Collection Collection
}
type Dependencies struct {
	Repositories repo.Repositories
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Collection: collection.NewCollectionService(deps.Repositories.Collection),
	}
}
