package rx

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func ToMap[K comparable, V any](it Iter[V], keyFn func(V) K) (m map[K]V, keys []K, err error) {
	m = make(map[K]V)
	err = ForEach(it, func(v V) error {
		k := keyFn(v)
		m[k] = v
		keys = append(keys, k)
		return nil
	})
	return
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
