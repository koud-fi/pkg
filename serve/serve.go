package serve

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/koud-fi/pkg/blob"
)

var ErrContentLengthMismatch = errors.New("content-length mismatch")

type Option func(c *config)

type config struct{ Info }

type Info struct {
	StatusCode    int
	ContentLength int64
	ContentType   string
	Location      string
	LastModified  time.Time

	// TODO: cache rules
}

func StatusCode(n int) Option         { return func(c *config) { c.StatusCode = n } }
func ContentLength(n int64) Option    { return func(c *config) { c.ContentLength = n } }
func ContentType(ct string) Option    { return func(c *config) { c.ContentType = ct } }
func Location(loc string) Option      { return func(c *config) { c.Location = loc } }
func LastModified(t time.Time) Option { return func(c *config) { c.LastModified = t } }

func Header(w http.ResponseWriter, opt ...Option) (*Info, error) {
	c := buildConfig(opt)
	c.writeHeader(w)
	return &c.Info, nil
}

func JSON(w http.ResponseWriter, r *http.Request, v any, opt ...Option) (*Info, error) {
	return Blob(w, r, blob.Marshal(json.Marshal, v),
		ContentType("application/json; charset=utf-8"))
}

func Blob(w http.ResponseWriter, r *http.Request, b blob.Blob, opt ...Option) (*Info, error) {
	rc, err := b.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var (
		c  = buildConfig(opt)
		rd io.Reader
	)
	if br, ok := rc.(blob.BytesReader); ok {
		buf := br.Bytes()
		c.ContentLength = int64(len(buf))

		// TODO: detect content-type

		rd = bytes.NewReader(buf)

	} else {
		if f, ok := rc.(fs.File); ok {
			info, err := f.Stat()
			if err != nil {
				return nil, err
			}
			c.ContentLength = info.Size()
			c.LastModified = info.ModTime()
		}

		// TODO: detect content-type

		rd = rc
	}
	c.writeHeader(w)

	// TODO: range requests

	if n, err := io.Copy(w, rd); err != nil {
		return nil, err
	} else if c.ContentLength > 0 && c.ContentLength != n {
		return nil, ErrContentLengthMismatch
	} else {
		c.ContentLength = n
	}
	return &c.Info, nil
}

func buildConfig(opt []Option) (c config) {
	for _, opt := range opt {
		opt(&c)
	}
	if c.StatusCode == 0 {
		c.StatusCode = http.StatusOK
	}
	return
}

func (nfo Info) writeHeader(w http.ResponseWriter) {
	if nfo.ContentLength > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(nfo.ContentLength, 10))
	}
	if nfo.ContentType != "" {
		w.Header().Set("Content-Type", nfo.ContentType)
	}
	if nfo.Location != "" {
		w.Header().Set("Location", nfo.Location)
	}
	if !nfo.LastModified.IsZero() {
		w.Header().Set("Last-Modified", nfo.LastModified.Format(http.TimeFormat))
	}
	w.WriteHeader(nfo.StatusCode)
}
