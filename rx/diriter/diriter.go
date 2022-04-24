package diriter

import (
	"io/fs"

	"github.com/koud-fi/pkg/rx"
)

func New(fsys fs.FS, root string) rx.Iter[string] {
	return rx.FuncIter(func() ([]string, bool, error) {

		// TODO: proper lazy implementation

		var paths []string
		err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		return paths, false, err
	})
}
