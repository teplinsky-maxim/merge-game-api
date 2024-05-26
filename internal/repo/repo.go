package repo

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo/repos"
	"merge-api/pkg/database"
)

type Collection interface {
	GetCollections(ctx context.Context, offset, limit uint) ([]entity.Collection, error)
	GetCollection(ctx context.Context, collectionId uint) (entity.CollectionWithItems, error)
	CreateCollection(ctx context.Context, name string) (entity.Collection, error)
	CreateCollectionItems(ctx context.Context, items []entity.CollectionItem) ([]entity.CollectionItem, error)
}

type Repositories struct {
	Collection
}

func NewRepositories(database *database.Database) *Repositories {
	return &Repositories{
		Collection: repos.NewCollectionRepo(database),
	}
}
