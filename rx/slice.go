package rx

type Slicer[T any] interface {
	Slice() []T
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
func (it sliceIter[T]) Slice() []T   { return it.data }

func Slice[T any](it Iter[T]) ([]T, error) {
	if s, ok := it.(Slicer[T]); ok {
		return s.Slice(), nil
	}
	return Reduce(it, func(s []T, v T) ([]T, error) { return append(s, v), nil }, []T{})
}

func UseSlice[T any](it Iter[T], fn func(s []T) error) error {
	s, err := Slice(it)
	if err != nil {
		return fn(s)
	}
	return err
}

/*
func Tee[T any](it Iter[T], out *[]T) Iter[T] {
	return Map(it, func(v T) T {
		*out = append(*out, v)
		return v
	})
}
*/
