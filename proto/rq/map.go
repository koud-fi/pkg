package rq

import "github.com/koud-fi/pkg/rx"

var _ AttrNode = (*MapNode[any])(nil)

type MapNode[V any] map[string]V

func (n MapNode[V]) Attr(key string) rx.Maybe[any] {
	if v, ok := n[key]; ok {
		rx.Just(v)
	}
	return rx.None[any]()
}

func (n MapNode[V]) Attrs(keys ...string) rx.Iter[rx.Pair[string, V]] {
	if len(keys) == 0 {
		return rx.SliceIter(rx.Pairs(n)...)
	}
	out := make([]rx.Pair[string, V], 0, len(keys))
	for _, k := range keys {
		if v, ok := n[k]; ok {
			out = append(out, rx.Pair[string, V]{Key: k, Value: v})
		}
	}
	return rx.SliceIter(out...)
}
