package rq

import "github.com/koud-fi/pkg/rx"

type Node interface {
	Attrs(keys ...string) rx.Iter[rx.Pair[string, any]]
}

type AttrNode interface {
	Attr(key string) rx.Maybe[any]
}

func Attr(n Node, key string) rx.Maybe[any] {
	if an, ok := n.(AttrNode); ok {
		return an.Attr(key)
	}
	v, _ := rx.First(n.Attrs(key))
	if !v.Ok() {
		return rx.None[any]()
	}
	return rx.Just(v.Value().Value)
}
