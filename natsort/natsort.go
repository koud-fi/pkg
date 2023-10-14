package natsort

import (
	"sort"
	"strconv"
)

func Strings(ss []string) {
	srts := make([]sortable, 0, len(ss))
	for _, s := range ss {
		srts = append(srts, toSortable(s))
	}
	sort.Slice(srts, func(i, j int) bool {
		return srts[i].compare(srts[j])
	})
	for i, srt := range srts {
		ss[i] = srt.s
	}
}

type chunk struct {
	s     string
	isNum bool
	n     int
}

type sortable struct {
	s      string
	chunks []chunk
}

func toSortable(s string) sortable {
	var (
		srt   = sortable{s: s}
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

func (s sortable) compare(to sortable) bool {
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
