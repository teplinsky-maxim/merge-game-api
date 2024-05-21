package collection

import (
	"context"
	"merge-api/internal/repo"
)

type CollectionService struct {
	repo repo.Collection
}

func (c *CollectionService) GetCollection(ctx context.Context, offset, limit uint) {
	//TODO implement me
	panic("implement me")
}

func NewCollectionService(repo repo.Collection) *CollectionService {
	return &CollectionService{repo: repo}
}
