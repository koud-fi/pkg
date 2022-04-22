package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

func SliceIter[T any](data ...T) Iter[T] {
	return &sliceIter[T]{data: data, offset: -1}
}

type sliceIter[T any] struct {
	data   []T
	offset int
}

func (it *sliceIter[_]) Next() bool {
	it.offset++
	return it.offset < len(it.data)
}

func (it sliceIter[T]) Value() T     { return it.data[it.offset] }
func (it sliceIter[_]) Close() error { return nil }

func FuncIter[T any](fn func() ([]T, bool, error)) Iter[T] {
	return &funcIter[T]{fn: fn}
}

type funcIter[T any] struct {
	fn func() ([]T, bool, error)
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
	hasMore := true
	for hasMore && it.offset >= len(it.data) && it.err == nil {
		it.data, hasMore, it.err = it.fn()
		it.offset = 0
	}
	return hasMore && it.offset < len(it.data) && it.err == nil
}

func (it funcIter[_]) Close() error { return it.err }

func Counter[T Number](start, step T) Iter[T] {
	return FuncIter(func() ([]T, bool, error) {
		next := start
		start += step
		return []T{next}, true, nil
	})
}
