package rx

import "log"

func Transform[T1, T2 any](it Iter[T1], fn func(T1) ([]T2, bool, error)) Iter[T2] {
	return FuncIter(func() ([]T2, bool, error) {
		if !it.Next() {
			return nil, false, nil
		}
		return fn(it.Value())
	})
}

func Map[T1, T2 any](it Iter[T1], fn func(T1) T2) Iter[T2] {
	return Transform(it, func(v T1) ([]T2, bool, error) {
		return []T2{fn(v)}, true, nil
	})
}

func Filter[T any](it Iter[T], fn func(T) bool) Iter[T] {
	return Transform(it, func(v T) ([]T, bool, error) {
		if fn(v) {
			return []T{v}, true, nil
		}
		return nil, true, nil
	})
}

func Unique[T comparable](it Iter[T]) Iter[T] {
	set := make(map[T]struct{})
	return Filter(it, func(v T) bool {
		if _, ok := set[v]; ok {
			return false
		}
		set[v] = struct{}{}
		return true
	})

}

func Skip[T any](it Iter[T], n int) Iter[T] {
	return Transform(it, func(v T) ([]T, bool, error) {
		n--
		if n >= 0 {
			return nil, true, nil
		}
		return []T{v}, true, nil
	})
}

func Take[T any](it Iter[T], n int) Iter[T] {
	return Transform(it, func(v T) ([]T, bool, error) {
		n--
		return []T{v}, n >= 0, nil
	})
}

func Log[T any](it Iter[T], prefix string) Iter[T] {
	return Transform(it, func(v T) ([]T, bool, error) {
		if prefix == "" {
			log.Print(v)
		} else {
			log.Println(prefix, v)
		}
		return []T{v}, true, nil
	})
}