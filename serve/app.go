package serve

import (
	"errors"
	"io/fs"
	"mime"
	"net/http"
	"os"
	pathpkg "path"

	"github.com/koud-fi/pkg/blob"
)

const DefaultIndexFile = "index.html"

type appConfig struct {
	index string
}

type AppOption func(*appConfig)

func Index(name string) AppOption { return func(c *appConfig) { c.index = name } }

func App(fsys fs.FS, opt ...AppOption) http.Handler {
	c := appConfig{
		index: DefaultIndexFile,
	}
	for _, opt := range opt {
		opt(&c)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Handle(w, r, func() (*Info, error) {
			if r.Method != http.MethodGet {
				return nil, os.ErrInvalid
			}
			path := r.URL.Path
			info, err := fs.Stat(fsys, path)
			if err != nil {
				if !os.IsNotExist(errors.Unwrap(err)) {
					return nil, err
				}
			}
			switch {
			case info != nil && !info.IsDir():
				break
			default:
				path = c.index
			}
			var opts []Option
			if ext := pathpkg.Ext(path); ext != "" {
				opts = append(opts, ContentType(mime.TypeByExtension(ext)))
			}
			return Blob(w, r, blob.FromFS(fsys, path), opts...)
		})
	})
}
