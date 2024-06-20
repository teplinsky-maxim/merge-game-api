package pkg

import "fmt"

type CollectionItem interface {
	GetCollectionInfo() (uint, uint)
	SetCollectionInfo(uint, uint)
	Equal(other CollectionItem) bool
	ToRedisString() string
}

type CollectionItemImpl struct {
	collectionId, collectionItemId uint
}

func (r *CollectionItemImpl) GetCollectionInfo() (uint, uint) {
	return r.collectionId, r.collectionItemId
}

func (r *CollectionItemImpl) SetCollectionInfo(collectionId, collectionItemId uint) {
	r.collectionId = collectionId
	r.collectionItemId = collectionItemId
}

func (r *CollectionItemImpl) Equal(other CollectionItem) bool {
	collectionId, collectionItemId := other.GetCollectionInfo()
	return r.collectionId == collectionId && r.collectionItemId == collectionItemId
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
