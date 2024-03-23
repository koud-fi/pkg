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

// HasSubset return true if b is a subset of a.
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
		if b[i] == a[j] {
			i++
		} else if b[i] < a[j] {
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

func (s Set[T]) Data() []T { return s.data }

// TODO: common set operations
