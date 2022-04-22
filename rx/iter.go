package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

func From[T any](data ...T) Iter[T] {
	return &sliceIter[T]{data: data, offset: -1}
}

type sliceIter[T any] struct {
	data   []T
	offset int
}

func (it *sliceIter[T]) Next() bool {
	it.offset++
	return it.offset < len(it.data)
}

func (it sliceIter[T]) Value() T     { return it.data[it.offset] }
func (it sliceIter[T]) Close() error { return nil }
