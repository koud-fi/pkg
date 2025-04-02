package logger

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	lrw.size += size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

type Logger struct {
	logFn LogFunc
}

func New() *Logger {
	return &Logger{
		logFn: DefaultLogFunc,
	}
}

func (l Logger) WithLogFunc(logFn LogFunc) *Logger {
	l.logFn = logFn
	return &l
}

func (l *Logger) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			start = time.Now()
			lrw   = &loggingResponseWriter{w, http.StatusOK, 0}
		)
		next.ServeHTTP(lrw, r)

		l.logFn(RequestInfo{
			Request:          r,
			Duration:         time.Since(start),
			ResponseStatus:   lrw.statusCode,
			ResponseBodySize: lrw.size,
		})
	})
}

type LogFunc func(RequestInfo)

func DefaultLogFunc(ri RequestInfo) {

	// TODO: color coding for errors

	log.Printf("[%d] %s %s %s %d %v",
		ri.ResponseStatus, ri.Request.RemoteAddr, ri.Request.Method,
		ri.PrettyURL(), ri.ResponseBodySize, ri.Duration)
}

// TODO: create a log func that uses slog, with extension possibilities

type RequestInfo struct {
	Request          *http.Request
	Duration         time.Duration
	ResponseStatus   int
	ResponseBodySize int
}

func (ri RequestInfo) PrettyURL() (url string) {
	if ri.Request.Host != "localhost" {
		url = ri.Request.Host
	}
	url += ri.Request.URL.Path
	if ri.Request.URL.RawQuery != "" {
		url += "?" + ri.Request.URL.RawQuery
	}
	return url
}
