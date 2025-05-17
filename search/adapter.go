package search

import "hash/crc64"

type Adapter[T any] struct {
	idFn    func(v T) string
	tagsFn  func(v T) []string
	orderFn func(v T) int64
}

func NewAdapter[T any](idFn func(v T) string, tagsFn func(v T) []string) Adapter[T] {
	return Adapter[T]{
		idFn:    idFn,
		tagsFn:  tagsFn,
		orderFn: defaultOrderFn(idFn),
	}
}

func (a Adapter[T]) WithOrderFn(orderFn func(v T) int64) Adapter[T] {
	a.orderFn = orderFn
	return a
}

func (a Adapter[T]) ID(v T) string     { return a.idFn(v) }
func (a Adapter[T]) Tags(v T) []string { return a.tagsFn(v) }
func (a Adapter[T]) Order(v T) int64   { return a.orderFn(v) }

func defaultOrderFn[T any](idFn func(v T) string) func(v T) int64 {
	return func(v T) int64 {
		id := idFn(v)
		h := crc64.Checksum([]byte(id), shardHashKeyTable)
		return int64(h)
	}
}
