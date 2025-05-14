package fetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/koud-fi/pkg/blob"
)

type Middleware func(*http.Request) (*http.Request, error)

func (r *Request) append(m Middleware) *Request {
	r.middlewares = append(r.middlewares, m)
	return r
}

func setUrl(u string) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		if r.URL.String() != "" { // TODO: this requires more robust handling

			// TODO: check if u is absolute or relative, and act accordingly

			u = r.URL.String() + u
		} else {
			switch {
			case strings.HasPrefix(u, "//"):
				u = "http:" + u
			case strings.HasPrefix(u, "/"):
				u = "http://localhost" + u
			}
		}
		if r.URL, err = url.Parse(u); err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		return r, nil
	}
}

func setContext(ctx context.Context) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		return r.WithContext(ctx), nil
	}
}

func setMethod(m string) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		r.Method = m
		return r, nil
	}
}

func setQuery(key string, value any) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {

		// TODO: this is not the best way to handle this, does allocation every time

		q := r.URL.Query()
		q.Set(key, fmt.Sprint(value))
		r.URL.RawQuery = q.Encode()
		return r, nil
	}
}

func setHeader(key string, value any) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		setReqHeader(r, key, fmt.Sprint(value))
		return r, nil
	}
}

func setUser(u *url.Userinfo) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		pass, passSet := u.Password()
		// Prevent setting invalid credentials.
		if u.Username() != "" && (!passSet || pass != "") {
			r.URL.User = u
		}
		return r, nil
	}
}

func setBody(b blob.Reader, mime string) Middleware {
	return func(r *http.Request) (_ *http.Request, err error) {
		data, err := blob.Bytes(b)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}

		// TODO: avoid reading the entire body to memory to resolve content length

		r.Body = io.NopCloser(bytes.NewReader(data))

		setReqHeader(r, "Content-Type", mime)
		setReqHeader(r, "Content-Length", strconv.Itoa(len(data)))
		return r, nil
	}
}

func formReader(data url.Values) blob.Reader {
	return blob.FromString(data.Encode())
}

func setReqHeader(r *http.Request, key, value string) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set(key, value)
}
