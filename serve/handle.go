package serve

import (
	"log"
	"net/http"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request, fn func() (*Info, error)) {
	startTime := time.Now()

	// TODO: improve size logging (both incoming and outgoing)
	// TODO: make all this shit configurable

	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: proper CORS configs

	info, err := fn()
	if info == nil && err != nil {
		Error(w, err)
		log.Printf("[%d] %s %s (%v) ERROR: %v",
			ErrorStatusCode(err), r.Method, r.URL, time.Since(startTime), err)
	} else {
		log.Printf("[%d] %s %s (%d bytes, %v)",
			info.StatusCode, r.Method, r.URL, info.ContentLength, time.Since(startTime))
	}
}
