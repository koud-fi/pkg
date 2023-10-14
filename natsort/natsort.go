// Package natsort implements somewhat optimized "natural sorting" functions
// https://en.wikipedia.org/wiki/Natural_sort_order
package natsort

import (
	"sort"
	"strconv"
)

// Strings sorts a list of string to natural order
func Strings(ss []string) {
	Values[string](ss, func(s string) string { return s })
}

// Values sorts a list of arbitary values to natural order,
// strFn is used to determine the string used for sorting for each value
func Values[T any](vs []T, strFn func(v T) string) {
	srts := make([]sortable[T], 0, len(vs))
	for _, v := range vs {
		srts = append(srts, toSortable(v, strFn(v)))
	}
	sort.Slice(srts, func(i, j int) bool {
		return srts[i].compare(srts[j])
	})
	for i, srt := range srts {
		vs[i] = srt.v
	}
}

type chunk struct {
	s     string
	isNum bool
	n     int
}

type sortable[T any] struct {
	v      T
	s      string
	chunks []chunk
}

func toSortable[T any](v T, s string) sortable[T] {
	var (
		srt   = sortable[T]{v: v, s: s}
		start int
		isNum bool
	)
	for i := 0; i <= len(s); i++ {
		if i == 0 {
			isNum = isNumber(s[i])
		} else {
			nextIsNum := i < len(s) && isNumber(s[i])
			if nextIsNum != isNum || i == len(s) {
				c := chunk{s: s[start:i]}
				if isNum {
					n, err := strconv.Atoi(c.s)
					c.n, c.isNum = n, err == nil
				}
				srt.chunks = append(srt.chunks, c)
				isNum = nextIsNum
				start = i
			}
		}
	}
	return srt
}

func isNumber(b byte) bool {
	return b >= 48 && b <= 57
}

func (s sortable[T]) compare(to sortable[T]) bool {
	for i := range s.chunks {
		if i >= len(to.chunks) {
			return false
		}
		a, b := s.chunks[i], to.chunks[i]

		if a.isNum && b.isNum {
			if a.n == b.n {
				switch {
				case i == len(s.chunks)-1:
					return true
				case i == len(to.chunks)-1:
					return false
				default:
					continue
				}
			}
			return a.n < b.n
		}
		if a.s == b.s {
			switch {
			case i == len(s.chunks)-1:
				return true
			case i == len(to.chunks)-1:
				return false
			default:
				continue
			}
		}
		return a.s < b.s

	}
	return false
}
