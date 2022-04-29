package num

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Gompertz function based on https://en.wikipedia.org/wiki/Gompertz_function
func Gompertz(t, a, b, c float64) float64 {
	return a * math.Pow(math.E, b*math.Pow(math.E, c*t))
}

func Clamp[T constraints.Ordered](n, min, max T) T {
	if n < min {
		return min
	} else if n > max {
		return max
	}
	return n
}
