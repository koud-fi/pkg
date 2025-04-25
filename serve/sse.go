package serve

import (
	"errors"
	"fmt"
	"iter"
	"net/http"

	"github.com/koud-fi/pkg/blob"
)

type SSEEvent struct {
	Event string
	Data  blob.Reader
}

func SSE[T any](
	w http.ResponseWriter,
	r *http.Request,
	iter iter.Seq[T],
	mapFn func(T) SSEEvent,
) (*Info, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	nfo := &Info{
		ContentType: "text/event-stream",
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	nfo.writeHeader(w)

	for v := range iter {
		e := mapFn(v)
		data, err := blob.String(e.Data)
		if err != nil {
			return nil, fmt.Errorf("read data: %w", err)
		}
		w.Write([]byte("event: " + e.Event + "\n"))
		w.Write([]byte("data: " + data + "\n\n"))
		flusher.Flush()
	}
	return nfo, nil
}
