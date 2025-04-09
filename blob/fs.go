package blob

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func FromFS(fsys fs.FS, name string) Reader {
	return Func(func() (io.ReadCloser, error) {
		return fsys.Open(name)
	})
}

func FromFile(name string) Reader {
	return Func(func() (io.ReadCloser, error) {
		absPath, err := filepath.Abs(name)
		if err != nil {
			return nil, err
		}
		return os.Open(absPath)
	})
}
