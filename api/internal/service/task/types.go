package task

type CreateNewBoardTaskInput struct {
	Width  uint `json:"width,omitempty"`
	Height uint `json:"height,omitempty"`
}

type MoveItemTaskInput struct {
	BoardID uint `json:"board_id,omitempty"`
	W1      uint `json:"w1,omitempty"`
	H1      uint `json:"h1,omitempty"`
	W2      uint `json:"w2,omitempty"`
	H2      uint `json:"h2,omitempty"`
}

type MergeItemsTaskInput struct {
	BoardID uint `json:"board_id,omitempty"`
	W1      uint `json:"w1,omitempty"`
	H1      uint `json:"h1,omitempty"`
	W2      uint `json:"w2,omitempty"`
	H2      uint `json:"h2,omitempty"`
}
