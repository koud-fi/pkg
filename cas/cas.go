package cas

import (
	"bytes"
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/grf"
)

type Storage struct {
	s        blob.Storage
	g        *grf.Graph
	nt       grf.NodeType
	fileOpts []file.Option
}

func New(s blob.Storage, g *grf.Graph, nt grf.NodeType, fileOps ...file.Option) *Storage {
	return &Storage{s: s, g: g, nt: nt, fileOpts: fileOps}
}

func (s *Storage) Lookup(id ID) (*Node, error) {
	d, err := s.g.MappedNode(s.nt, id.String(), false).Data()
	if err != nil {
		return nil, err
	}
	nd := d.(NodeData)
	if nd.Size == 0 {
		return nil, grf.ErrNotFound
	}
	return &Node{ID: id, NodeData: nd, s: s.s}, nil
}

func (s *Storage) Add(b blob.Blob) (*Node, error) {

	// TODO: avoid full memory copy

	data, err := blob.Bytes(b)
	if err != nil {
		return nil, err
	}
	var (
		id  = NewIDFromBytes(data)
		key = id.String()
	)
	d, err := s.g.MappedNode(s.nt, key, true).Update(func(_ any) (any, error) {
		if err := s.s.Set(context.Background(), id.Hex(), bytes.NewReader(data)); err != nil {
			return nil, err
		}
		return file.ResolveAttrs(blob.FromBytes(data), s.fileOpts...)
	}).Data()
	if err != nil {
		return nil, err
	}
	return &Node{ID: id, NodeData: d.(NodeData), s: s.s}, nil
}
