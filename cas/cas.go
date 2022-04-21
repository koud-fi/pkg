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
	n, err := grf.Mapped[NodeData](s.g, s.nt, id.String())
	if err != nil {
		return nil, err
	}
	if n.Data.Size == 0 {
		return nil, grf.ErrNotFound
	}
	return &Node{ID: id, NodeData: n.Data, s: s.s}, nil
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
	n, err := grf.SetMapped(s.g, s.nt, key, func(_ NodeData) (NodeData, error) {
		if err := s.s.Set(context.Background(), id.Hex(), bytes.NewReader(data)); err != nil {
			return NodeData{}, err
		}
		attrs, err := file.ResolveAttrs(blob.FromBytes(data), s.fileOpts...)
		if err != nil {
			return NodeData{}, err
		}
		return NodeData{Attributes: *attrs}, nil
	})
	if err != nil {
		return nil, err
	}
	return &Node{ID: id, NodeData: n.Data, s: s.s}, nil
}
