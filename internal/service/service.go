package service

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo"
	"merge-api/internal/service/collection"
)

type Collection interface {
	GetCollection(ctx context.Context, input *collection.GetCollectionInput) ([]entity.Collection, error)
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
