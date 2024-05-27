package board

type SizeType uint

type Board[T any] interface {
	Get(width, height SizeType) (T, error)
	Set(width, height SizeType, t T) error
	Width() SizeType
	Height() SizeType
}
