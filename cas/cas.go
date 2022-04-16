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
	return &Storage{s: s, g: g, nt: nt}
}

func (s *Storage) Lookup(id ID) (*Node, error) {
	n, err := s.g.MappedNode(s.nt, id.String())
	if err != nil {
		return nil, err
	}
	var attrs file.Attributes
	if err := n.Unmarshal(&attrs); err != nil {
		return nil, err
	}
	return &Node{ID: id, Attributes: attrs, s: s.s}, nil
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
	if n, err := s.g.MappedNode(s.nt, key); err == nil {
		var attrs file.Attributes
		if err := n.Unmarshal(&attrs); err != nil {
			return nil, err
		}
		return &Node{ID: id, Attributes: attrs, s: s.s}, nil

	} else if err != grf.ErrNotFound {
		return nil, err
	}
	attrs, err := file.ResolveAttrs(blob.FromBytes(data), s.fileOpts...)
	if err != nil {
		return nil, err
	}
	if err := s.s.Set(context.Background(), id.Hex(), bytes.NewReader(data)); err != nil {
		return nil, err
	}
	if _, err := s.g.AddMappedNode(s.nt, key, attrs); err != nil {
		return nil, err
	}
	return &Node{ID: id, Attributes: *attrs, s: s.s}, nil
}
