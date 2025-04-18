package timecamp

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/calendar"
	"github.com/koud-fi/pkg/fetch"
)

const apiRoot = "https://app.timecamp.com/third_party/api"

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

type Client struct {
	baseReq *fetch.Request
}

func New(apiKey string) *Client {
	return &Client{baseReq: fetch.New().Authorization(apiKey)}
}

func (c Client) CreateEntries(e ...Entry) error {
	for _, e := range e {
		if e.EndTime == nil {
			return errors.New("missing end time")
		}
		if err := blob.Error(c.baseReq.
			Method(http.MethodPost).
			URL(apiRoot + "/entries").
			Form(url.Values{
				"task_id":    []string{strconv.FormatInt(e.TaskID, 10)},
				"date":       []string{e.StartTime.Format("2006-01-02")},
				"start_time": []string{e.StartTime.Format(time.RFC3339)},
				"end_time":   []string{e.EndTime.Format(time.RFC3339)},
				"note":       []string{e.Note},
			}),
		); err != nil {
			return err
		}
	}

	// TODO: return IDs of the created entries?

	return nil
}
