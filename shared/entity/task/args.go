package task

import (
	"encoding/json"
)

type NewBoardTaskArgs struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

type NewBoardTaskResult struct {
	BoardID uint `json:"board_id"`
}

func (a *NewBoardTaskArgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}

func NewNewBoardTaskArgs(width, height uint) NewBoardTaskArgs {
	return NewBoardTaskArgs{
		Width:  width,
		Height: height,
	}
}

type MoveItemTaskArgs struct {
	BoardID uint `json:"board_id"`
	W1      uint `json:"w1"`
	H1      uint `json:"h1"`
	W2      uint `json:"w2"`
	H2      uint `json:"h2"`
}

type MoveItemTaskResult struct {
	Result bool   `json:"result"`
	Reason string `json:"reason"`
}

func (a *MoveItemTaskArgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}

func NewMoveItemTaskArgs(boardId, w1, h1, w2, h2 uint) MoveItemTaskArgs {
	return MoveItemTaskArgs{
		BoardID: boardId,
		W1:      w1,
		H1:      h1,
		W2:      w2,
		H2:      h2,
	}
}
