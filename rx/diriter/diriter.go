package diriter

import (
	"io/fs"

	"github.com/koud-fi/pkg/rx"
)

type Node struct {
	Path string
	fs.DirEntry
}

func New(fsys fs.FS, root string) rx.Iter[Node] {
	return rx.FuncIter(func() ([]Node, bool, error) {

		// TODO: proper lazy implementation

		var nodes []Node
		err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			nodes = append(nodes, Node{
				Path:     path,
				DirEntry: d,
			})
			return nil
		})
		return nodes, false, err
	})
}

func Paths(it rx.Iter[Node]) rx.Iter[string] {
	return rx.Map(it, func(n Node) string { return n.Path })
}
