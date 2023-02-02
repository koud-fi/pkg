package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

type Doner interface{ Done() bool }

type Done bool

func (d Done) Done() bool { return bool(d) }

func FuncIter[T any, S Doner](fn func() ([]T, S, error)) Iter[T] {
	return &funcIter[T, S]{fn: fn}
}

type funcIter[T any, S Doner] struct {
	fn    func() ([]T, S, error)
	state S
	sIter sliceIter[T]
	err   error
}

func (it *funcIter[_, _]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.sIter.Next() {
		return true
	}
	for len(it.sIter.data) == 0 && !it.state.Done() && it.err == nil {
		it.sIter.data, it.state, it.err = it.fn()
	}
	return (len(it.sIter.data) > 0 || !it.state.Done()) && it.err == nil
}

func (it funcIter[T, _]) Value() T     { return it.sIter.Value() }
func (it funcIter[_, _]) Close() error { return it.err }

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
	return FuncIter(func() ([]T, Done, error) {
		return nil, true, err
	})
}

func Counter[N Number](start, step N) Iter[N] {
	return FuncIter(func() ([]N, Done, error) {
		next := start
		start += step
		return []N{next}, false, nil
	})
}

func Range[N Number](start, step N, count int) Iter[N] {
	return Take(Counter(start, step), count)
}
