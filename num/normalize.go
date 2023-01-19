package num

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

// Gompertz function based on https://en.wikipedia.org/wiki/Gompertz_function
func Gompertz(t, a, b, c float64) float64 {
	return a * math.Pow(math.E, b*math.Pow(math.E, c*t))
}

// Nudge adds random variance by a factor (0.0 for not variance, 1.0 for completely random)
// to the given value, also calmping it between min and max.
func Nudge(v, min, max, factor float64, seed int64) float64 {
	if max < min {
		min, max = max, min
	}
	nudge := rand.New(rand.NewSource(seed)).Float64() * factor * (max - min)
	v += nudge - nudge*0.5
	return math.Min(max, math.Max(min, v))*0.99 + nudge*0.01
}

func Clamp[T constraints.Ordered](n, min, max T) T {
	if n < min {
		return min
	} else if n > max {
		return max
	}
	return n
}
