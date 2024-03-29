package rx

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

type Doner interface{ Done() bool }

type Done bool

func (d Done) Done() bool { return bool(d) }

type Forever struct{}

func (f Forever) Done() bool { return true }

type funcIter[T any, S Doner] struct {
	fn    func(S) ([]T, S, error)
	state S
	sIter sliceIter[T]
	err   error

	lensInit bool
	lens     Lens[S]
}

func FuncIter[T any, S Doner](fn func(S) ([]T, S, error)) Iter[T] {
	return &funcIter[T, S]{fn: fn}
}

func Unfold[T any, S Doner](s S, fn func(S) ([]T, S, error)) Iter[T] {
	return &funcIter[T, S]{fn: fn, state: s}
}

func UnfoldLens[T any, S Doner](l Lens[S], fn func(S) ([]T, S, error)) Iter[T] {
	return &funcIter[T, S]{fn: fn, lens: l}
}

func (it *funcIter[_, _]) Next() bool {
	if it.err != nil {
		return false
	}
	if it.sIter.Next() {
		return true
	}
	if it.lens != nil && !it.lensInit {
		if it.state, it.err = it.lens.Get(); it.err != nil {
			return false
		}
		it.lensInit = true
	}
	for len(it.sIter.data) == 0 && !it.state.Done() && it.err == nil {
		if it.lens != nil {
			if it.err = it.lens.Set(it.state); it.err != nil {
				return false
			}
		}
		it.sIter.data, it.state, it.err = it.fn(it.state)
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
	return FuncIter(func(Done) ([]T, Done, error) {
		return nil, true, err
	})
}

func Counter[N Number](start, step N) Iter[N] {
	return FuncIter(func(Done) ([]N, Done, error) {
		next := start
		start += step
		return []N{next}, false, nil
	})
}

func Range[N Number](start, step N, count int) Iter[N] {
	return Take(Counter(start, step), count)
}
