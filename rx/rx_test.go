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

func TestFuncIter(t *testing.T) {
	rx.Drain(rx.Log(rx.FuncIter(func() ([]int, bool, error) {
		return []int{1, 2, 3}, false, nil
	}), ""))
}

func TestUnique(t *testing.T) {
	t.Log(rx.Slice(rx.Unique(rx.SliceIter(5, 2, 1, 2, 3, 3, 1, 4, 1))))
}

func TestDistinct(t *testing.T) {
	t.Log(rx.Slice(rx.Distinct(rx.SliceIter(0, 0, 2, 1, 1, 1, 2, 3, 2, 2, 3, 3, 4))))
}

func TestSkipAndTake(t *testing.T) {
	t.Log(rx.Slice(rx.Take(rx.Skip(rx.Counter(1, 1), 10), 10)))
}

func TestParitionAll(t *testing.T) {
	t.Log(rx.Slice(rx.PartitionAll(rx.Range(1, 1, 10), 3)))
}

func TestParitionLoops(t *testing.T) {
	t.Log(rx.Slice(rx.PartitionLoops(rx.SliceIter(0, 1, 2, 3, 4, 2, 3, 4, 0, 5, 4, 1))))
}

func TestSum(t *testing.T) {
	t.Log(rx.Sum(rx.Range(1, 1, 10)))
}

func TestToMapUtils(t *testing.T) {
	pairs := rx.Pluck(rx.Range(1, 1, 5), func(n int) string {
		return "." + strconv.Itoa(n)
	})
	m, _ := rx.ToMap(pairs)
	pairSlice := rx.SortedPairs(m, rx.SortKeys[string, int])
	t.Log(pairSlice)
	/*
		vs, _ := rx.Slice(rx.SelectKeys(m, rx.SliceIter(keys...)))
		t.Log(m, keys, vs)
	*/
}
