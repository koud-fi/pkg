package rq

import "github.com/koud-fi/pkg/rx"

type Node interface {
	Attrs(keys ...string) rx.Iter[rx.Pair[string, any]]
}

type AttrNode interface {
	Attr(key string) any
}

func Attr(n Node, key string) any {
	v, _, err := rx.First(n.Attrs(key))
	if err != nil {
		return err
	}
	return v
}
