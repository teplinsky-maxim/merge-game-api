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
	ID           uint
	CollectionId uint
	Name         string
	Level        uint8
	Mergeable    bool // means item can be merged with another item of the same item
	CanCreate    bool // means item can create another items using click
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
