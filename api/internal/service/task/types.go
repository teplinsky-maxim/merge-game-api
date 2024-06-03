package task

import (
	"merge-api/shared/pkg/board"
)

type CreateNewBoardTaskInput struct {
	Width  board.SizeType
	Height board.SizeType
}
