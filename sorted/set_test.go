package sorted_test

import (
	"testing"

	"github.com/koud-fi/pkg/sorted"
)

func TestHasSubset(t *testing.T) {
	assert(t, "1", sorted.NewSet[int]().HasSubset(sorted.NewSet[int]()), true)

	assert(t, "2", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet[int]()), true)
	assert(t, "3", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(1)), true)
	assert(t, "4", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(2, 4)), true)
	assert(t, "5", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(0, 4)), false)
	assert(t, "6", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(2, 6)), false)
	assert(t, "7", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(0)), false)
	assert(t, "8", sorted.NewSet(1, 2, 3, 4, 5).HasSubset(sorted.NewSet(7)), false)

}

func assert(t *testing.T, tag string, result, expected bool) {
	if result != expected {
		t.Fatalf("%s: result mismatch", tag)
	}
}
