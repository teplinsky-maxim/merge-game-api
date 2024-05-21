package collection

import (
	"context"
	"merge-api/internal/entity"
	"merge-api/internal/repo"
)

type CollectionService struct {
	repo repo.Collection
}

func (c *CollectionService) GetCollection(ctx context.Context, input *GetCollectionInput) ([]entity.Collection, error) {
	if input.Limit == nil {
		limit := uint(100)
		input.Limit = &limit
	}
	if input.Offset == nil {
		offset := uint(0)
		input.Offset = &offset
	}
	collection, err := c.repo.GetCollection(ctx, *input.Offset, *input.Limit)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func NewCollectionService(repo repo.Collection) *CollectionService {
	return &CollectionService{repo: repo}
}
