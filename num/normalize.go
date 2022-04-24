package num

import "math"

// Gompertz function based on https://en.wikipedia.org/wiki/Gompertz_function
func Gompertz(t, a, b, c float64) float64 {
	return a * math.Pow(math.E, b*math.Pow(math.E, c*t))
}
