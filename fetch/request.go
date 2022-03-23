package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/koud-fi/pkg/blob"

	"golang.org/x/time/rate"
)

var _ blob.Blob = (*Request)(nil)

type Request struct {
	client   *http.Client
	ctx      context.Context
	method   string
	url      string
	query    []pair
	header   []pair
	user     *url.Userinfo
	limiter  *rate.Limiter
	body     blob.Blob
	bodyMime string
}

type pair struct {
	key    string
	values []interface{}
}

func New() *Request {
	return &Request{client: http.DefaultClient}
}

func Get(url string) *Request    { return New().Method(http.MethodGet).URL(url) }
func Head(url string) *Request   { return New().Method(http.MethodHead).URL(url) }
func Post(url string) *Request   { return New().Method(http.MethodPost).URL(url) }
func Put(url string) *Request    { return New().Method(http.MethodPut).URL(url) }
func Delete(url string) *Request { return New().Method(http.MethodDelete).URL(url) }

func (r Request) Method(m string) *Request {
	r.method = m
	return &r
}

func (r Request) URL(u string) *Request {
	switch {
	case strings.HasPrefix(u, "//"):
		r.url = "http:" + u
	case strings.HasPrefix(u, "/"):
		r.url = "http://localhost" + u
	default:
		r.url = u
	}
	return &r
}

func (r Request) Client(c *http.Client) *Request {
	r.client = c
	return &r
}

func (r Request) Context(ctx context.Context) *Request {
	r.ctx = ctx
	return &r
}

func (r Request) Query(key string, vs ...interface{}) *Request {
	r.query = append(r.query, pair{key: key, values: vs})
	return &r
}

func (r Request) Header(key string, vs ...interface{}) *Request {
	r.header = append(r.header, pair{key: http.CanonicalHeaderKey(key), values: vs})
	return &r
}

func (r Request) UserAgent(ua string) *Request {
	return r.Header("User-Agent", ua)
}

func (r Request) Authorization(authHeader string) *Request {
	return r.Header("Authorization", authHeader)
}

func (r Request) User(u *url.Userinfo) *Request {
	pass, passSet := u.Password()
	if u.Username() != "" && (!passSet || pass != "") {
		r.user = u
	}
	return &r
}

func (r Request) Body(b blob.Blob, mime string) *Request {
	r.body = b
	r.bodyMime = mime
	return &r
}

func (r Request) Form(data url.Values) *Request {
	r.body = blob.FromString(data.Encode())
	r.bodyMime = "application/x-www-form-urlencoded; charset=utf-8"
	return &r
}

func (r Request) Limit(l *rate.Limiter) *Request {
	r.limiter = l
	return &r
}

func (r *Request) HttpRequest() (*http.Request, error) {
	u, err := url.Parse(strings.SplitN(r.url, "?", 2)[0])
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if rawQuery := injectPairs(u.Query(), r.query).Encode(); rawQuery != "" {
		u.RawQuery = rawQuery
	}
	u.User = r.user

	var body io.Reader
	if r.body != nil {
		var err error
		if body, err = r.body.Open(); err != nil {
			return nil, fmt.Errorf("failed to open body: %w", err)
		}
	}
	if r.method == "" {
		r.method = http.MethodGet
	}
	req, err := http.NewRequest(r.method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}
	injectPairs(url.Values(req.Header), r.header)

	if r.bodyMime != "" {
		req.Header.Set("Content-Type", r.bodyMime)
	}
	return req, nil
}

func injectPairs(vals url.Values, ps []pair) url.Values {
	for _, p := range ps {
		for _, vi := range p.values {
			var vStr string
			switch v := vi.(type) {
			case []byte:
				vStr = string(v)
			default:
				vStr = fmt.Sprint(v)
			}
			if vStr == "" {
				continue
			}
			vals[p.key] = append(vals[p.key], vStr)
		}
	}
	return vals
}
