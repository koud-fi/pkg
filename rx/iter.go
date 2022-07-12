package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

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

func FuncIter[T any](fn func() ([]T, bool, error)) Iter[T] {
	return &funcIter[T]{fn: fn, hasMore: true}
}

type funcIter[T any] struct {
	fn func() ([]T, bool, error)
	sliceIter[T]
	hasMore bool
	err     error
}

func (it *funcIter[_]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.sliceIter.Next() {
		return true
	}
	for len(it.data) == 0 && it.hasMore && it.err == nil {
		it.data, it.hasMore, it.err = it.fn()
	}
	return (len(it.data) > 0 || it.hasMore) && it.err == nil
}

func (it funcIter[_]) Close() error { return it.err }

func WithClose[T any](it Iter[T], fn func() error) Iter[T] {
	return closeIter[T]{Iter: it, closeFn: fn}
}

type closeIter[T any] struct {
	Iter[T]
	closeFn func() error
}

func (it closeIter[_]) Close() error {
	if err := it.Iter.Close(); err != nil {
		it.closeFn()
		return err
	}
	return it.closeFn()
}

func Counter[N Number](start, step N) Iter[N] {
	return FuncIter(func() ([]N, bool, error) {
		next := start
		start += step
		return []N{next}, true, nil
	})
}

func Range[N Number](start, step N, count int) Iter[N] {
	return Take(Counter(start, step), count)
}
