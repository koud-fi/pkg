package diriter

import (
	"io/fs"
	"path"

	"github.com/koud-fi/pkg/rx"
)

func New(fsys fs.FS, root string) rx.Iter[string] {
	return rx.FuncIter(func() ([]string, bool, error) {
		dir, err := fs.ReadDir(fsys, root)
		if err != nil {
			return nil, false, err
		}
		out := make([]string, 0, len(dir))
		for _, d := range dir {
			out = append(out, path.Join(root, d.Name()))
		}
		return out, false, nil
	})
}
