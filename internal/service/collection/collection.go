package collection

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo"
)

type CollectionService struct {
	repo repo.Collection
}

func (c *CollectionService) CreateCollection(ctx context.Context, input *CreateCollectionInput) (entity.Collection, error) {
	collection, err := c.repo.CreateCollection(ctx, (*input).Name)
	if err != nil {
		return entity.Collection{}, nil
	}
	return collection, nil
}

func (c *CollectionService) GetCollections(ctx context.Context, input *GetCollectionsInput) ([]entity.Collection, error) {
	if input.Limit == nil {
		limit := uint(100)
		input.Limit = &limit
	}
	if input.Offset == nil {
		offset := uint(0)
		input.Offset = &offset
	}
	collection, err := c.repo.GetCollections(ctx, *input.Offset, *input.Limit)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (c *CollectionService) GetCollection(ctx context.Context, input *GetCollectionInput) (entity.CollectionWithItems, error) {
	collection, err := c.repo.GetCollection(ctx, (*input).Id)
	if err != nil {
		return entity.CollectionWithItems{}, nil
	}
	return collection, nil
}

func NewCollectionService(repo repo.Collection) *CollectionService {
	return &CollectionService{repo: repo}
}
