package rx

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Random[N constraints.Integer](min, max N) Iter[N] {
	return FuncIter(func() ([]N, bool, error) {
		return []N{N(rand.Int63n(int64(max)-int64(min))) + min}, true, nil
	})
}
