package rx

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