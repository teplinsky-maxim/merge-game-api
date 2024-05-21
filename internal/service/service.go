package service

import (
	"context"
	"merge-api/internal/repo"
	"merge-api/internal/service/collection"
)

type Collection interface {
	GetCollection(ctx context.Context, offset, limit uint)
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
