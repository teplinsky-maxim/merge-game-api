package collection

type GetCollectionInput struct {
	Limit  *uint `json:"limit,omitempty"`
	Offset *uint `json:"offset,omitempty"`
}
