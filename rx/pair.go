package rx

import "golang.org/x/exp/constraints"

type Pair[K, V any] struct {
	key   K
	value V
}

func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{key, value}
}

func (p Pair[K, _]) Key() K   { return p.key }
func (p Pair[_, V]) Value() V { return p.value }

func Pluck[K, V any](it Iter[V], keyFn func(V) K) Iter[Pair[K, V]] {
	return PluckErr(it, func(v V) (K, error) { return keyFn(v), nil })
}

func PluckErr[K, V any](it Iter[V], keyFn func(V) (K, error)) Iter[Pair[K, V]] {
	return MapErr(it, func(v V) (p Pair[K, V], err error) {
		p.value = v
		p.key, err = keyFn(v)
		return
	})
}

func SortKeys[K constraints.Ordered, V any](a, b Pair[K, V]) bool {
	return a.key < b.key
}

func SortKeysDesc[K constraints.Ordered, V any](a, b Pair[K, V]) bool {
	return a.key > b.key
}

func SortValues[K comparable, V constraints.Ordered](a, b Pair[K, V]) bool {
	return a.value < b.value
}

func SortValuesDesc[K comparable, V constraints.Ordered](a, b Pair[K, V]) bool {
	return a.value > b.value
}
