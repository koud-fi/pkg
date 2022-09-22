package router

import "strings"

type Router[T any] struct {
	// TODO
}

func (r Router[T]) Lookup(path ...string) (T, func(string) string, bool) {
	frags := make([]string, 0, len(path))
	for _, p := range path {
		frags = append(frags, strings.Split(p, "/")...)
	}

	// ???

	panic("TODO")
}
