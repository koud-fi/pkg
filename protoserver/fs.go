package protoserver

import (
	"context"
	"io"
	"io/fs"
	pathpkg "path"
	"strings"

	"github.com/koud-fi/pkg/pk"
	"github.com/koud-fi/pkg/schema"
)

func RegisterFS(s pk.Scheme, fsys fs.FS) {
	Register(s, FetchFunc(func(_ context.Context, ref pk.Ref) (any, error) {
		return &FSNode{fsys: fsys, scheme: ref.Scheme(), path: ref.Key()}, nil
	}))
}

var _ interface {
	schema.RefNode
	schema.ParentNode
	schema.FileNode
} = (*FSNode)(nil)

type FSNode struct {
	fsys   fs.FS
	scheme pk.Scheme
	path   string
}

func (n FSNode) Ref() (pk.Ref, error) {
	return pk.NewRef(n.scheme, "", n.path)
}

func (n FSNode) Children() ([]schema.RefNode, error) {
	dir, err := fs.ReadDir(n.fsys, n.path)
	if err != nil {
		return nil, err
	}
	nodes := make([]schema.RefNode, 0, len(dir))
	for _, d := range dir {
		if defaultHideFunc(d.Name()) {
			continue
		}
		nodes = append(nodes, FSNode{
			fsys:   n.fsys,
			scheme: n.scheme,
			path:   pathpkg.Join(n.path, d.Name()),
		})
	}
	return nodes, nil
}

func defaultHideFunc(name string) bool {
	return strings.HasPrefix(name, ".")
}

func (n FSNode) Files() (map[string]schema.File, error) {
	return map[string]schema.File{schema.MasterFileKey: n}, nil
}

func (n FSNode) Open() (io.ReadCloser, error) {
	return n.fsys.Open(n.path)
}
