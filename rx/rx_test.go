package rx_test

import (
	"strconv"
	"testing"

	"github.com/koud-fi/pkg/rx"
)

func TestIter(t *testing.T) {
	n := rx.Counter(0, 1)
	n = rx.Skip(n, 10)
	n = rx.Log(n, "1.")
	n = rx.Filter(n, func(n int) bool { return n%2 == 0 })
	s := rx.Map(n, func(n int) string { return "N" + strconv.Itoa(n*2) })
	s = rx.Take(s, 5)
	rx.Drain(rx.Log(s, "2."))
}

func TestSum(t *testing.T) {
	t.Log(rx.Sum(rx.Take(rx.Counter(1, 1), 10)))
}
