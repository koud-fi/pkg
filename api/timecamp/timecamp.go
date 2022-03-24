package timecamp

import (
	"time"

	"github.com/koud-fi/pkg/calendar"
)

var _ calendar.MutablePeriod = (*Entry)(nil)

type Entry struct {
	TaskID    int64
	StartTime time.Time
	EndTime   *time.Time
	Note      string
}

func (e Entry) Start() time.Time      { return e.StartTime }
func (e Entry) End() *time.Time       { return e.EndTime }
func (e *Entry) SetStart(t time.Time) { e.StartTime = t }
func (e *Entry) SetEnd(t *time.Time)  { e.EndTime = t }

// TODO
