package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

func FuncIter[T any](fn func() ([]T, bool, error)) Iter[T] {
	return &funcIter[T]{fn: fn, hasMore: true}
}

type funcIter[T any] struct {
	fn      func() ([]T, bool, error)
	sIter   sliceIter[T]
	hasMore bool
	err     error
}

func (it *funcIter[_]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.sIter.Next() {
		return true
	}
	for len(it.sIter.data) == 0 && it.hasMore && it.err == nil {
		it.sIter.data, it.hasMore, it.err = it.fn()
	}
	return (len(it.sIter.data) > 0 || it.hasMore) && it.err == nil
}

func (it funcIter[T]) Value() T     { return it.sIter.Value() }
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

func Error[T any](err error) Iter[T] {
	return FuncIter(func() ([]T, bool, error) {
		return nil, false, err
	})
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
