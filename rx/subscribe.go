package rx

func First[T any](it Iter[T]) (Maybe[T], error) {
	if !it.Next() {
		return None[T](), it.Close()
	}
	return Just(it.Value()), nil
}

func ForEach[T any](it Iter[T], fn func(T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}

func ForEachN[T any](it Iter[T], fn func(T, int) error) error {
	i := -1
	return ForEach(it, func(v T) error {
		i++
		return fn(v, i)
	})
}

func Reduce[T, S any](it Iter[T], fn func(S, T) (S, error), sum S) (S, error) {
	err := ForEach(it, func(v T) (err error) {
		sum, err = fn(sum, v)
		return
	})
	return sum, err
}

func Sum[N Number](it Iter[N]) (N, error) {
	return Reduce(it, func(sum, n N) (N, error) { return sum + n, nil }, 0)
}

func Count[T any](it Iter[T]) (int, error) {
	var n int
	err := ForEach(it, func(T) error { n++; return nil })
	return n, err
}

func Drain[T any](it Iter[T]) error {
	return ForEach(it, func(T) error { return nil })
}
