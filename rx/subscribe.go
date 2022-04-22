package rx

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}

func Reduce[T, S any](it Iter[T], fn func(S, T) (S, error), sum S) (S, error) {
	err := ForEach(it, func(v T) (err error) {
		sum, err = fn(sum, v)
		return
	})
	return sum, err
}

func Slice[T any](it Iter[T]) ([]T, error) {
	return Reduce(it, func(s []T, v T) ([]T, error) { return append(s, v), nil }, []T{})
}

func Sum[T Number](it Iter[T]) (T, error) {
	return Reduce(it, func(sum, n T) (T, error) { return sum + n, nil }, 0)
}

func Drain[T any](it Iter[T]) {
	ForEach(it, func(_ T) error { return nil })
}
