package rx

func Partition[T any](it Iter[T], fn func([]T) ([]T, [][]T)) Iter[[]T] {
	var vs []T
	return FuncIter(func(Done) ([][]T, Done, error) {
		for it.Next() {
			var out [][]T
			if vs, out = fn(append(vs, it.Value())); len(out) > 0 {
				return out, false, nil
			}
		}
		if len(vs) > 0 {
			return [][]T{vs}, true, it.Close()
		}
		return nil, true, it.Close()

	})
}

func PartitionAll[T any](it Iter[T], n int) Iter[[]T] {
	return Partition(it, func(vs []T) ([]T, [][]T) {
		if len(vs) == n {
			return nil, [][]T{vs}
		}
		return vs, nil
	})
}

func PartitionLoops[T comparable](it Iter[T]) Iter[[]T] {
	return Partition(it, func(vs []T) ([]T, [][]T) {
		last := len(vs) - 1
		for i := last - 1; i >= 0; i-- {
			if vs[i] == vs[last] {
				if i > 0 {
					return []T{vs[last]}, [][]T{vs[:i], vs[i:last]}
				}
				return nil, [][]T{vs}
			}
		}
		return vs, nil
	})
}
