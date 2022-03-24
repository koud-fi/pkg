package calendar

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	StartTime time.Time
	EndTime   *time.Time
	Body      string
}

func (e Entry) Start() time.Time      { return e.StartTime }
func (e Entry) End() *time.Time       { return e.EndTime }
func (e *Entry) SetStart(t time.Time) { e.StartTime = t }
func (e *Entry) SetEnd(t *time.Time)  { e.EndTime = t }

func Parse(s string, loc *time.Location) (e Entry, err error) {
	parts := strings.Fields(s)
	if parts, err = e.parsePeriod(parts, loc); err != nil {
		return
	}
	e.Body = strings.Join(parts, " ")
	return
}

func (e *Entry) parsePeriod(s []string, loc *time.Location) ([]string, error) {
	if len(s) < 2 {
		return nil, errors.New("invalid period")
	}
	start, err := parseTime(s[0]+" "+s[1], loc)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}
	e.SetStart(start)
	if len(s) > 2 {
		if end, err := parseTime(s[0]+" "+s[2], loc); err == nil {
			e.SetEnd(&end)
			s = s[1:]
		}
	}
	return s[2:], nil
}

func parseTime(s string, loc *time.Location) (time.Time, error) {
	if t, err := time.Parse("2006-1-2 15", s); err == nil {
		return t, nil
	}
	return time.Parse("2006-1-2 15:04", s)
}

func (e Entry) String() string {
	var sb strings.Builder
	sb.WriteString(e.Start().Format("2006-1-2") + " " + formatTime(e.Start()))
	if e.End() != nil {
		sb.WriteString(formatTime(*e.End()))
	}
	if e.Body != "" {
		sb.WriteByte(' ')
		sb.WriteString(e.Body)
	}
	return sb.String()
}

func formatTime(t time.Time) string {
	if t.Minute() == 0 {
		return t.Format("15")
	}
	return t.Format("15:04")
}
