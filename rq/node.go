package rq

import "github.com/koud-fi/pkg/rx"

type Node interface {
	Attrs(keys ...string) rx.Iter[rx.Pair[string, any]]
}

type AttrNode interface {
	Attr(key string) (any, bool)
}

func Attr(n Node, key string) (any, bool) {
	p, ok, err := rx.First(n.Attrs(key))
	if err != nil {
		return err, false
	}
	return p.Value, ok
}
