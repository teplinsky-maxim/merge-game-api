package collection

import "merge-api/api/internal/entity"

type GetCollectionsInput struct {
	Limit  *uint `json:"limit,omitempty"`
	Offset *uint `json:"offset,omitempty"`
}

type GetCollectionInput struct {
	Id uint
}

type CreateCollectionInput struct {
	Name string
}

type CreateCollectionItemsInput struct {
	Items []entity.CollectionItem
}
