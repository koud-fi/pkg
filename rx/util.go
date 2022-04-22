package rx

func ForEach[T any](it Iter[T], fn func(v T) error) error {
	for it.Next() {
		if err := fn(it.Value()); err != nil {
			return err
		}
	}
	return it.Close()
}
