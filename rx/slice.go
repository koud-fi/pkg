package rx

func SliceIter[T any](data ...T) Iter[T] {
	return &sliceIter[T]{data: data}
}

type sliceIter[T any] struct {
	data []T
	init bool
}

func (it *sliceIter[_]) Next() bool {
	if !it.init {
		it.init = true
	} else {
		it.data = it.data[1:]
	}
	return len(it.data) > 0
}

func (it sliceIter[T]) Value() T     { return it.data[0] }
func (it sliceIter[_]) Close() error { return nil }

func Slice[T any](it Iter[T]) ([]T, error) {
	return Reduce(it, func(s []T, v T) ([]T, error) { return append(s, v), nil }, []T{})
}
