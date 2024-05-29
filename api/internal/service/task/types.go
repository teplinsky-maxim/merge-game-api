package task

import "merge-api/pkg/board"

type CreateNewBoardTaskInput struct {
	Width  board.SizeType
	Height board.SizeType
}
