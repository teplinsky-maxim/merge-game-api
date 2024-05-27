package task

import (
	"encoding/json"
)

type Args interface {
	MarshalJSON() ([]byte, error)
}
type NewBoardTaskArgs struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

func (a *NewBoardTaskArgs) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func NewNewBoardTaskArgs(width, height uint) NewBoardTaskArgs {
	return NewBoardTaskArgs{
		Width:  width,
		Height: height,
	}
}
