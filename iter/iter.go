package iter

type Iter[T any] interface {
	Next() bool
	Value() T
	Close() error
}

func From[T any](data ...T) Iter[T] {
	return &dataIter[T]{data: data, offset: -1}
}

type dataIter[T any] struct {
	data   []T
	offset int
}

func (it *dataIter[T]) Next() bool {
	it.offset++
	return it.offset < len(it.data)
}

func (it dataIter[T]) Value() T     { return it.data[it.offset] }
func (it dataIter[T]) Close() error { return nil }

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}
