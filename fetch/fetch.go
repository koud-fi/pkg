package fetch // import "github.com/hyde042/private/lib/go/core/fetch"

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/koud-fi/pkg/blob"
	"golang.org/x/time/rate"
)

const (
	jsonMime = "application/json"
	formMime = "application/x-www-form-urlencoded; charset=utf-8"
)

var _ blob.Reader = (*Request)(nil)

func Get(url string) *Request    { return New().Get(url) }
func Head(url string) *Request   { return New().Head(url) }
func Post(url string) *Request   { return New().Post(url) }
func Put(url string) *Request    { return New().Put(url) }
func Patch(url string) *Request  { return New().Patch(url) }
func Delete(url string) *Request { return New().Delete(url) }

// URI will automatically resolve the correct fetching mechanism and return a blob.Reader for the underlying data.
func URI(uri string) blob.Reader { return defaultURIFetcher.Fetch(uri) }

func (r Request) Get(url string) *Request    { return r.Method(http.MethodGet).URL(url) }
func (r Request) Head(url string) *Request   { return r.Method(http.MethodHead).URL(url) }
func (r Request) Post(url string) *Request   { return r.Method(http.MethodPost).URL(url) }
func (r Request) Put(url string) *Request    { return r.Method(http.MethodPut).URL(url) }
func (r Request) Patch(url string) *Request  { return r.Method(http.MethodPatch).URL(url) }
func (r Request) Delete(url string) *Request { return r.Method(http.MethodDelete).URL(url) }

func (r Request) Client(c *http.Client) *Request       { return r.setClient(c) }
func (r Request) Context(ctx context.Context) *Request { return r.append(setContext(ctx)) }
func (r Request) Limit(l *rate.Limiter) *Request       { return r.setLimiter(l) }
func (r Request) Middleware(m Middleware) *Request     { return r.append(m) }

// TODO: custom error handlers (a post-process system?)
// TODO: support for custom cookie jars? (what would be the benefit? is this correct layer?)

func (r Request) Method(m string) *Request         { return r.append(setMethod(m)) }
func (r Request) URL(u string) *Request            { return r.append(setUrl(u)) }
func (r Request) Query(key string, v any) *Request { return r.append(setQuery(key, v)) }

func (r Request) Header(key string, v any) *Request        { return r.append(setHeader(key, v)) }
func (r Request) UserAgent(ua string) *Request             { return r.Header("User-Agent", ua) }
func (r Request) Authorization(header string) *Request     { return r.Header("Authorization", header) }
func (r Request) BasicAuth(user, password string) *Request { return r.setBasicAuth(user, password) }
func (r Request) BearerAuth(token string) *Request         { return r.Authorization("Bearer " + token) }

func (r Request) User(u *url.Userinfo) *Request { return r.append(setUser(u)) }

func (r Request) Body(b blob.Reader, mime string) *Request { return r.append(setBody(b, mime)) }
func (r Request) JSON(v any) *Request                      { return r.append(setBody(blob.MarshalJSON(v), jsonMime)) }
func (r Request) Form(data url.Values) *Request            { return r.append(setBody(formReader(data), formMime)) }

// TODO: multi-part form

func (r Request) DirReader(dr DirReader) *Request { return r.setDirReader(dr) }

func (r *Request) HTTPRequest() (*http.Request, error)        { return r.httpRequest() }
func (r *Request) Do() (*http.Response, *http.Request, error) { return r.do() }

func (r *Request) Open() (io.ReadCloser, error) { return r.openFile() } // blob.Reader
func (r *Request) OpenFile() (fs.File, error)   { return r.openFile() }
func (r *Request) Stat() (fs.FileInfo, error)   { return r.stat() }
