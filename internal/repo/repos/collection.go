package repos

import (
	"context"
	"merge-api/pkg/database"
)

type CollectionRepo struct {
	database *database.Database
}

func (c *CollectionRepo) GetCollection(ctx context.Context, offset, limit uint) {}

func NewCollectionRepo(database *database.Database) *CollectionRepo {
	return &CollectionRepo{
		database: database,
	}
}
