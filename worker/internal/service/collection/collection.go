package collection

import (
	"context"
	"merge-api/worker/internal/repo"
	"merge-api/worker/pkg"
)

type CollectionService struct {
	repo *repo.Collection
}

func (r *CollectionService) GetNextCollectionItem(ctx context.Context, item pkg.CollectionItem) (pkg.CollectionItem, error) {
	return (*r.repo).GetNextCollectionItem(ctx, item)
}

func (r *CollectionService) IsItemMergeable(ctx context.Context, item pkg.CollectionItem) (bool, error) {
	return (*r.repo).IsItemMergeable(ctx, item)
}

func NewCollectionService(repo repo.Collection) *CollectionService {
	return &CollectionService{
		repo: &repo,
	}
}
