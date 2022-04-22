package rx

func Transform[T1, T2 any](it Iter[T1], fn func(T1) ([]T2, error)) Iter[T2] {
	return FuncIter(func() ([]T2, error) {
		if !it.Next() {
			return nil, nil
		}
		return fn(it.Value())
	})
}

func Take[T any](it Iter[T], n int) Iter[T] {
	return Transform(it, func(v T) ([]T, error) {
		if n == 0 {
			return nil, nil
		}
		n--
		return []T{v}, nil
	})
}
