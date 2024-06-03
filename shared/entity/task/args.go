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
