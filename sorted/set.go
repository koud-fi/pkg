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

	return HasSubset(a.data, b.data)
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

func (s *Set[T]) Data() []T {
	s.clean()
	return s.data
}

// HasSubset return true if b is a subset of a.
// a and b MUST BE sorted in ascending order.
func HasSubset[T constraints.Ordered](a, b []T) bool {
	var (
		al = len(a)
		bl = len(b)
	)
	if al == 0 {
		return bl == 0
	}
	if bl == 0 {
		return true
	}
	var i, j int
	if b[bl-1] > a[al-1] {
		return false
	}
	for i < bl {

		// TODO: use binary search to find the next value?

		if b[i] == a[j] {
			i++
		} else if b[i] < a[j] {
			return false
		}
		j++
	}
	return true
}

// TODO: common set operations
