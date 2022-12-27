package rq

import "github.com/koud-fi/pkg/rx"

type Node interface {
	Tag() string
	Attr(key string) any
	Children(tag ...string) rx.Iter[Node]
}
