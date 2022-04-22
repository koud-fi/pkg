package rx_test

import (
	"strconv"
	"testing"

	"github.com/koud-fi/pkg/rx"
)

func TestIter(t *testing.T) {
	n := rx.Counter(1, 1)
	s := rx.Map(n, func(n int) string { return "N" + strconv.Itoa(n*2) })
	s = rx.Take(s, 10)
	rx.Log(s)
}
