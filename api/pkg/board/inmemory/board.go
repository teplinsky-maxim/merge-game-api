package inmemory

import "errors"

var BoardIndexError = errors.New("index out of board")

type Board[T any] struct {
	width  uint
	height uint
	grid   [][]T
}

// coordsAreInBounds is to check if oob error
func coordsAreInBounds[T any](width, height uint, grid *[][]T) bool {
	if height >= uint(len(*grid)) {
		return false
	}
	tmp := (*grid)[height]
	return width < uint(len(tmp))
}

func (b *Board[T]) Get(width, height uint) (T, error) {
	if coordsAreInBounds[T](width, height, &b.grid) {
		return b.grid[height][width], nil
	}
	return *new(T), BoardIndexError
}

func (b *Board[T]) Set(width, height uint, t T) error {
	if !coordsAreInBounds[T](width, height, &b.grid) {
		return BoardIndexError
	}
	b.grid[height][width] = t
	return nil
}

func (b *Board[T]) Width() uint {
	return b.width
}

func (b *Board[T]) Height() uint {
	return b.height
}

func NewBoard[T any](width, height uint) Board[T] {
	grid := make([][]T, height)
	for i := range grid {
		grid[i] = make([]T, width)
	}
	return Board[T]{
		width:  width,
		height: height,
		grid:   grid,
	}
}
