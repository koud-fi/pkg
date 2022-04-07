package bloom_test

import (
	"testing"

	"github.com/koud-fi/pkg/bloom"
)

func TestFilter(t *testing.T) {
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 1))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 2))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 3))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 10))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 100))

	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 9).Contains(bloom.New([]uint64{0}, 9)))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 9).Contains(bloom.New([]uint64{0, 64}, 9)))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 9).Contains(bloom.New([]uint64{100, 150}, 9)))
	t.Log(bloom.New([]uint64{0, 63, 64, 255, 257}, 9).Contains(bloom.New([]uint64{1}, 9)))
}
