package collection

type GetCollectionsInput struct {
	Limit  *uint `json:"limit,omitempty"`
	Offset *uint `json:"offset,omitempty"`
}

type GetCollectionInput struct {
	Id uint
}
