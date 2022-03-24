package calendar

import (
	"errors"
	"sort"
	"time"
)

type Period interface {
	Start() time.Time
	End() *time.Time
}

type MutablePeriod interface {
	Period
	SetStart(time.Time)
	SetEnd(*time.Time)
}

func Normalize(ps []MutablePeriod) error {

	// TODO: add options for time "smoothing" and overlap elmination

	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Start().Before(ps[j].Start())
	})
	for i := range ps {

		// TODO: add options to nil end if it's same as next periods start?

		if ps[i].End() == nil {
			next := i + 1
			if next > len(ps)-1 {
				return errors.New("missing end from last period")
			}
			nextStart := ps[next].Start()
			ps[i].SetEnd(&nextStart)
		}
	}
	return nil
}
