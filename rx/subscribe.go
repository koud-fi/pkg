package rx

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}

func Reduce[T, S any](it Iter[T], fn func(S, T) (S, error)) (s S, err error) {
	err = ForEach(it, func(v T) (err error) {
		s, err = fn(s, v)
		return
	})
	return
}

func Sum[T Number](it Iter[T]) (T, error) {
	return Reduce(it, func(sum, n T) (T, error) { return sum + n, nil })
}

func Drain[T any](it Iter[T]) {
	ForEach(it, func(_ T) error { return nil })
}
