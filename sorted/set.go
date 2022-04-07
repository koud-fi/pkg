package sorted

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Set[T constraints.Ordered] struct {
	data    []T
	isDirty bool
}

func NewSet[T constraints.Ordered](v ...T) *Set[T] {
	s := new(Set[T])
	s.Add(v...)
	return s
}

// Add does not check for member uniqueness.
func (s *Set[T]) Add(v ...T) {
	s.data = append(s.data, v...)
	s.isDirty = true
}

func (a *Set[T]) HasSubset(b *Set[T]) bool {
	a.clean()
	b.clean()
	var (
		al = len(a.data)
		bl = len(b.data)
	)
	if al == 0 {
		return bl == 0
	}
	if bl == 0 {
		return true
	}
	var i, j int
	if b.data[bl-1] > a.data[al-1] {
		return false
	}
	for i < bl {
		if b.data[i] == a.data[j] {
			i++
		} else if b.data[i] < a.data[j] {
			return false
		}
		j++
	}
	return true
}

func (s *Set[T]) clean() {
	if s.isDirty {

		// TODO: use "generic" sort?

		sort.Slice(s.data, func(i, j int) bool {
			return s.data[i] < s.data[j]
		})
		s.isDirty = false
	}
}

// TODO: common set operations
