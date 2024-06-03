package task

import "merge-api/api/pkg/board"

type CreateNewBoardTaskInput struct {
	Width  board.SizeType
	Height board.SizeType
}
