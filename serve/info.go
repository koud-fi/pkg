package serve

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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
	Range         *[2]int64
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
	if nfo.Range != nil {
		w.Header().Set("Range", fmt.Sprintf("%d-%d", nfo.Range[0], nfo.Range[1]))
	}
	w.WriteHeader(nfo.StatusCode)
}
