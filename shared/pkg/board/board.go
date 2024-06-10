package board

type Board[T any] interface {
	Get(width, height uint) (T, error)
	Set(width, height uint, t T) error
	Width() uint
	Height() uint
}
