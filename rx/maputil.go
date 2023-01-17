package rx

func ToMap[K comparable, V any](it Iter[Pair[K, V]]) (map[K]V, error) {
	m := make(map[K]V)
	return m, ForEach(it, func(p Pair[K, V]) error {
		m[p.key] = p.value
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
