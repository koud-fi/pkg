package fetch

import (
	"encoding/base64"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

type Request struct {
	client      *http.Client
	limiter     *rate.Limiter
	middlewares []Middleware
	dirReader   DirReader
}

func New() *Request {
	return &Request{client: http.DefaultClient}
}

func (r Request) setClient(c *http.Client) *Request {
	r.client = c
	return &r
}

func (r Request) setLimiter(l *rate.Limiter) *Request {
	r.limiter = l
	return &r
}

func (r Request) setBasicAuth(username, password string) *Request {
	header := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return r.Authorization("Basic " + header)
}

func (r Request) setDirReader(dr DirReader) *Request {
	r.dirReader = dr
	return &r
}

func (r *Request) httpRequest() (*http.Request, error) {
	req := &http.Request{
		URL: new(url.URL),
	}
	for _, m := range r.middlewares {
		var err error
		if req, err = m(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}
