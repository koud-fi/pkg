package rx

import "log"

func Transform[T1, T2 any](it Iter[T1], fn func(T1) ([]T2, Done, error)) Iter[T2] {
	return FuncIter(func(Done) ([]T2, Done, error) {
		if !it.Next() {
			return nil, true, it.Close()
		}
		return fn(it.Value())
	})
}

func Any[T any](it Iter[T]) Iter[any] {
	return Transform(it, func(v T) ([]any, Done, error) {
		return []any{v}, false, nil
	})
}

func Map[T1, T2 any](it Iter[T1], fn func(T1) T2) Iter[T2] {
	return MapErr(it, func(v T1) (T2, error) { return fn(v), nil })
}

func MapErr[T1, T2 any](it Iter[T1], fn func(T1) (T2, error)) Iter[T2] {
	return Transform(it, func(v T1) ([]T2, Done, error) {
		out, err := fn(v)
		if err != nil {
			return nil, true, err
		}
		return []T2{out}, false, nil
	})
}

func Filter[T any](it Iter[T], fn func(T) bool) Iter[T] {
	return Transform(it, func(v T) ([]T, Done, error) {
		if fn(v) {
			return []T{v}, false, nil
		}
		return nil, false, nil
	})
}

func Unique[T comparable](it Iter[T]) Iter[T] {
	return UniqueFunc(it, func(v T) T { return v })
}

func UniqueFunc[T any, K comparable](it Iter[T], fn func(T) K) Iter[T] {
	set := make(map[K]struct{})
	return Filter(it, func(v T) bool {
		k := fn(v)
		if _, ok := set[k]; ok {
			return false
		}
		set[k] = struct{}{}
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
	return Transform(it, func(v T) ([]T, Done, error) {
		n--
		return []T{v}, n <= 0, nil
	})
}

func Log[T any](it Iter[T], prefix string) Iter[T] {
	return Transform(it, func(v T) ([]T, Done, error) {
		if prefix == "" {
			log.Print(v)
		} else {
			log.Println(prefix, v)
		}
		return []T{v}, false, nil
	})
}
