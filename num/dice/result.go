package dice

import "github.com/koud-fi/pkg/num"

type Result struct {
	Die   Die
	Rolls []int
}

func (r Result) Total() int { return num.Sum(r.Rolls...) }
