package iter

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}
