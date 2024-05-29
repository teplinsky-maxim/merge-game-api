package entity

// Collection is a model
type Collection struct {
	ID   uint
	Name string
}

type CollectionWithItems struct {
	Collection
	Items []CollectionItemWithGenerate
}

// CollectionItem is a model
type CollectionItem struct {
	ID           uint   `json:"id,omitempty"`
	CollectionId uint   `json:"collection_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Level        uint8  `json:"level,omitempty"`
	Mergeable    bool   `json:"mergeable,omitempty"`  // means item can be merged with another item of the same item
	CanCreate    bool   `json:"can_create,omitempty"` // means item can create another items using click
}

type CollectionItemWithGenerate struct {
	CollectionItem
	Generate []struct {
		CollectionID uint
		Level        uint
	}
}

// CreationRule is a model
type CreationRule struct {
	ID                       uint
	CollectionItemId         uint
	GenerateCollectionItemId uint
}
