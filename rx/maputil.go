package rx

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func Pluck[K comparable, V any](it Iter[V], keyFn func(V) K) Iter[Pair[K, V]] {
	return PluckErr(it, func(v V) (K, error) { return keyFn(v), nil })
}

func PluckErr[K comparable, V any](it Iter[V], keyFn func(V) (K, error)) Iter[Pair[K, V]] {
	return MapErr(it, func(v V) (p Pair[K, V], err error) {
		p.Value = v
		p.Key, err = keyFn(v)
		return
	})
}

func ToMap[K comparable, V any](it Iter[Pair[K, V]]) (map[K]V, error) {
	m := make(map[K]V)
	return m, ForEach(it, func(p Pair[K, V]) error {
		m[p.Key] = p.Value
		return nil
	})
}

func SelectKeys[K comparable, V any](m map[K]V, keys Iter[K]) Iter[V] {
	return Map(keys, func(k K) V { return m[k] })
}

func Pairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, Pair[K, V]{k, v})
	}
	return pairs
}

type PairSortFunc[K comparable, V any] func(a, b Pair[K, V]) bool

func SortedPairs[K comparable, V any](m map[K]V, fn PairSortFunc[K, V]) []Pair[K, V] {
	pairs := Pairs(m)
	sort.Slice(pairs, func(i, j int) bool {
		return fn(pairs[i], pairs[j])
	})
	return pairs
}

func SortKeys[K constraints.Ordered, V any](a, b Pair[K, V]) bool {
	return a.Key < b.Key
}

func SortKeysDesc[K constraints.Ordered, V any](a, b Pair[K, V]) bool {
	return a.Key > b.Key
}

func SortValues[K comparable, V constraints.Ordered](a, b Pair[K, V]) bool {
	return a.Value < b.Value
}

func SortValuesDesc[K comparable, V constraints.Ordered](a, b Pair[K, V]) bool {
	return a.Value > b.Value
}

func Key[K comparable, V any](p Pair[K, V]) K   { return p.Key }
func Value[K comparable, V any](p Pair[K, V]) V { return p.Value }
