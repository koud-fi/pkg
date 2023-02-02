package rx

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Random[N constraints.Integer](min, max N) Iter[N] {
	max++
	return FuncIter(func(Done) ([]N, Done, error) {
		return []N{N(rand.Int63n(int64(max)-int64(min))) + min}, false, nil
	})
}
