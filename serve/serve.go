package serve

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/koud-fi/pkg/blob"
)

var ErrContentLengthMismatch = errors.New("content-length mismatch")

var rangeParser = regexp.MustCompile(`bytes=(\d+)-(\d+)?`)

type Option func(c *config)

type config struct{ Info }

type Info struct {
	StatusCode    int
	ContentLength int64
	ContentType   string
	Location      string
	LastModified  time.Time
	MaxAge        time.Duration
	Immutable     bool
	Compress      bool
	Disposition   string
}

func StatusCode(n int) Option         { return func(c *config) { c.StatusCode = n } }
func ContentLength(n int64) Option    { return func(c *config) { c.ContentLength = n } }
func ContentType(ct string) Option    { return func(c *config) { c.ContentType = ct } }
func Location(loc string) Option      { return func(c *config) { c.Location = loc } }
func LastModified(t time.Time) Option { return func(c *config) { c.LastModified = t } }
func MaxAge(d time.Duration) Option   { return func(c *config) { c.MaxAge = d } }
func Compress(b bool) Option          { return func(c *config) { c.Compress = b } }
func Immutable(b bool) Option         { return func(c *config) { c.Immutable = b } }

func Attachment(name string) Option {
	return func(c *config) {
		if name == "" {
			c.Disposition = ""
		}
		c.Disposition = fmt.Sprintf(`attachment; filename="%s"`, name)
	}
}

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
	return Reader(w, r, rc, opt...)
}

func Reader(rw http.ResponseWriter, r *http.Request, rd io.Reader, opt ...Option) (*Info, error) {
	var c = buildConfig(opt)
	if br, ok := rd.(blob.BytesReader); ok {
		buf := br.Bytes()
		c.ContentLength = int64(len(buf))

		// TODO: detect content-type

		rd = bytes.NewReader(buf)

	} else {
		if f, ok := rd.(fs.File); ok {
			info, err := f.Stat()
			if err != nil {
				return nil, err
			}
			c.ContentLength = info.Size()
			c.LastModified = info.ModTime()
		}

		// TODO: detect content-type

	}
	w := io.Writer(rw)

	// TODO: support range requests with io.Seeker

	if rat, ok := rd.(io.ReaderAt); ok {
		rw.Header().Set("Accept-Ranges", "bytes")

		if rh := r.Header.Get("Range"); rh != "" {
			begin, end := parseRangeHeader(rh, c.ContentLength)
			buf := make([]byte, end-begin)
			if _, err := rat.ReadAt(buf, begin); err != nil {
				return nil, err
			}
			rw.Header().Set("Content-Range",
				fmt.Sprintf("bytes %d-%d/%d", begin, end-1, c.ContentLength))
			c.StatusCode = http.StatusPartialContent

			w = bytes.NewBuffer(buf)
		}
	}
	if c.Compress {
		rw.Header().Set("Vary", "Content-Encoding")
		if c.ContentLength > 1400 && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			c.ContentLength = -1
			rw.Header().Set("Content-Encoding", "gzip")

			gw := gzip.NewWriter(w)
			defer gw.Flush()
			w = gw

		} else {
			c.Compress = false
		}
	}
	c.writeHeader(rw)

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
	if nfo.MaxAge > 0 || nfo.Immutable {
		if nfo.Immutable && nfo.MaxAge == 0 {
			nfo.MaxAge = time.Hour * 24 * 365 // year
		}
		v := "public, max-age=" + strconv.Itoa(int(nfo.MaxAge/time.Second))
		if nfo.Immutable {
			v += ", immutable"
		}
		w.Header().Set("Cache-Control", v)
	}
	if nfo.Disposition != "" {
		w.Header().Set("Content-Disposition", nfo.Disposition)
	}
	w.WriteHeader(nfo.StatusCode)
}

func parseRangeHeader(s string, size int64) (int64, int64) {
	const bufSize = 1 << 22 // 4 MB
	var (
		sms      = rangeParser.FindStringSubmatch(s)
		begin, _ = strconv.ParseInt(sms[1], 10, 64)
		end, _   = strconv.ParseInt(sms[2], 10, 64)
	)
	if end <= 0 || end == begin || end-begin > bufSize {
		end = begin + bufSize
		if end > size {
			end = size
		}
	}
	return begin, end
}
