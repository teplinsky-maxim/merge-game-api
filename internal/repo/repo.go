package repo

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo/repos"
	"merge-api/pkg/database"
)

type Collection interface {
	GetCollection(ctx context.Context, offset, limit uint) ([]entity.Collection, error)
}

type Repositories struct {
	Collection
}

func NewRepositories(database *database.Database) *Repositories {
	return &Repositories{
		Collection: repos.NewCollectionRepo(database),
	}
}
