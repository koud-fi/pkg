package rx

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}

func Discard[T any](it Iter[T]) {
	ForEach(it, func(_ T) error { return nil })
}
