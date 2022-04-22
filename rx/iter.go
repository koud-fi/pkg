package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

func SliceIter[T any](data ...T) Iter[T] {
	return &sliceIter[T]{data: data}
}

type sliceIter[T any] struct{ data []T }

func (it *sliceIter[_]) Next() bool {
	if len(it.data) == 0 {
		return false
	}
	it.data = (it.data)[:len(it.data)-1]
	return len(it.data) > 0
}

func (it sliceIter[T]) Value() T     { return it.data[0] }
func (it sliceIter[_]) Close() error { return nil }

func FuncIter[T any](fn func() ([]T, error)) Iter[T] {
	return &funcIter[T]{fn: fn}
}

type funcIter[T any] struct {
	fn func() ([]T, error)
	sliceIter[T]
	err error
}

func (it *funcIter[_]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.sliceIter.Next() {
		return true
	}
	it.data, it.err = it.fn()
	return len(it.data) > 0 && it.err == nil
}

func (it funcIter[_]) Close() error { return it.err }
