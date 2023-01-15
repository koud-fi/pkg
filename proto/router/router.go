package router

import (
	"strings"

	"github.com/koud-fi/pkg/rx"
)

const (
	PathSeparator = '/'
	ParamPrefix   = ':'
)

type Router[T any] struct {
	children map[string]Router[T]
	// TODO: isParam bool (or routeType?)
	value T
}

func (r *Router[T]) Register(path string, v T) {
	if r.children == nil {
		r.children = make(map[string]Router[T])
	}
	//part, rest := splitPath(path)

	// ???

	panic("TODO")
}

func (r *Router[T]) Lookup(path ...string) (T, Params, bool) {

	// ???

	panic("TODO")
}

func splitPath(path string) (string, string) {
	for len(path) > 0 && path[0] == PathSeparator {
		path = path[1:]
	}
	if idx := strings.IndexByte(path, PathSeparator); idx > -1 {
		return path[:idx], path[idx+1:]
	}
	return path, ""
}

type Params struct {
	pairs []rx.Pair[string, string]
}

func (p Params) Get(key string) string {
	for _, pair := range p.pairs {
		if pair.Key == key {
			return pair.Value
		}
	}
	panic("router: undefined parameter: " + key)
}
