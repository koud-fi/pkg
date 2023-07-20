package serve

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/koud-fi/pkg/blob"
)

const DefaultIndexFile = "index.html"

type appConfig struct {
	index string
}

type AppOption func(*appConfig)

func Index(name string) AppOption { return func(c *appConfig) { c.index = name } }

func App(opt ...AppOption) func(http.ResponseWriter, *http.Request, fs.FS, ...Option) (*Info, error) {
	c := appConfig{
		index: DefaultIndexFile,
	}
	for _, opt := range opt {
		opt(&c)
	}
	return func(w http.ResponseWriter, r *http.Request, fsys fs.FS, opt ...Option) (*Info, error) {
		if r.Method != http.MethodGet {
			return nil, os.ErrInvalid
		}
		p := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		info, err := fs.Stat(fsys, p)
		if err != nil {
			if !os.IsNotExist(errors.Unwrap(err)) {
				return nil, err
			}
		}
		switch {
		case info != nil && !info.IsDir():
			break
		default:
			p = c.index
		}
		opt = append(opt, ContentTypeFromPath(p))
		return Blob(w, r, blob.FromFS(fsys, p), opt...)
	}
}

func NewAppHandler(fsys fs.FS, opt ...AppOption) http.Handler {
	fn := App(opt...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r, func() (*Info, error) { return fn(w, r, fsys) })
	})
}
