package pkg

import "fmt"

type CollectionItem interface{}

type CollectionItemImpl struct {
	collectionId, collectionItemId uint
}

func (r *CollectionItemImpl) ToRedisString() string {
	return fmt.Sprintf("%v:%v", r.collectionId, r.collectionItemId)
}

func NewCollectionItemImpl(collectionId, collectionItemId uint) CollectionItemImpl {
	return CollectionItemImpl{
		collectionId:     collectionId,
		collectionItemId: collectionItemId,
	}
}
