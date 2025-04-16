package fetch

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/koud-fi/pkg/blob"
)

var defaultURIFetcher = new(URIFetcher)

type URIFetcher struct {
	fetchers map[string]URIFetchFunc
}

type URIFetchFunc func(string) (io.ReadCloser, error)

func (f URIFetcher) Register(schema string, fn URIFetchFunc) *URIFetcher {
	if f.fetchers == nil {
		f.fetchers = make(map[string]URIFetchFunc, 1)
	} else if _, ok := f.fetchers[schema]; ok {
		panic("duplicate registration for schema: " + schema)
	}
	f.fetchers[schema] = fn
	return &f
}

func (f *URIFetcher) Fetch(uri string) blob.Reader {
	return blob.Func(func() (io.ReadCloser, error) {
		switch {
		case strings.HasPrefix(uri, "//"):
			uri = "http:" + uri
		}
		parsedURI, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("invalid URI %q: %w", uri, err)
		}
		if f.fetchers != nil {
			if fetchFn, ok := f.fetchers[parsedURI.Scheme]; ok {
				return fetchFn(uri)
			}
		}
		switch parsedURI.Scheme {
		case "", "file":
			return blob.FromFile(strings.TrimPrefix(uri, "file://")).Open()
		case "http", "https":
			return Get(uri).Open()
		default:
			return nil, fmt.Errorf("unsupported schema %q", parsedURI.Scheme)
		}
	})
}

// TODO: fetcher for magnet links (requires quite complicated torrent client implementation)
