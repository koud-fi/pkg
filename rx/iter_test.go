package rx_test

import (
	"testing"

	"github.com/koud-fi/pkg/rx"
)

func TestIter(t *testing.T) {
	rx.Log(rx.Take(rx.Counter(1, 10), 10))
}
