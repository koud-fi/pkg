package simplex_test

import (
	"testing"

	"github.com/koud-fi/pkg/noise/simplex"
)

func TestNoise1D(t *testing.T) {
	for i := 0; i <= 20; i++ {
		t.Log(i, simplex.Noise1D(float32(i)/20))
	}
}
