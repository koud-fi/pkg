package rx

import "log"

func Transform[T1, T2 any](it Iter[T1], fn func(T1) ([]T2, bool, error)) Iter[T2] {
	return FuncIter(func() ([]T2, bool, error) {
		if !it.Next() {
			return nil, false, it.Close()
		}
		return fn(it.Value())
	})
}

func Any[T any](it Iter[T]) Iter[any] {
	return Transform(it, func(v T) ([]any, bool, error) {
		return []any{v}, true, nil
	})
}

func Map[T1, T2 any](it Iter[T1], fn func(T1) T2) Iter[T2] {
	return MapErr(it, func(v T1) (T2, error) { return fn(v), nil })
}

func MapErr[T1, T2 any](it Iter[T1], fn func(T1) (T2, error)) Iter[T2] {
	return Transform(it, func(v T1) ([]T2, bool, error) {
		out, err := fn(v)
		if err != nil {
			return nil, false, err
		}
		return []T2{out}, true, nil
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

func Distinct[T comparable](it Iter[T]) Iter[T] {
	return DistinctFunc(it, func(next, prev T) bool { return next != prev })
}

func DistinctFunc[T any](it Iter[T], fn func(T, T) bool) Iter[T] {
	var (
		init bool
		prev T
	)
	return Filter(it, func(v T) bool {
		if !init {
			init = true
		} else if !fn(v, prev) {
			return false
		}
		prev = v
		return true
	})
}

func Skip[T any](it Iter[T], n int) Iter[T] {
	return Filter(it, func(T) bool {
		n--
		return n < 0
	})
}

func Take[T any](it Iter[T], n int) Iter[T] {
	return Transform(it, func(v T) ([]T, bool, error) {
		n--
		return []T{v}, n > 0, nil
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
