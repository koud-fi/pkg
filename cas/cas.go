package cas

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/datastore"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/rx"
)

type Node struct {
	ID ID `json:"id"`
	file.Attributes
}

type Storage struct {
	s  blob.Storage
	ds *datastore.Sorted[file.Attributes]

	fileOpts []file.Option
}

func New(
	s blob.Storage,
	ds *datastore.Sorted[file.Attributes],
	fileOps ...file.Option,
) *Storage {
	return &Storage{s: s, ds: ds, fileOpts: fileOps}
}

func (s *Storage) Node(ctx context.Context, id ID) (*Node, error) {
	attrs, err := s.ds.Get(ctx, id.Hex())
	if err != nil {
		return nil, err
	}

	// TODO: ensure that the attributes are up to date

	return &Node{
		ID:         id,
		Attributes: attrs,
	}, nil
}

func (s *Storage) Get(ctx context.Context, id ID) blob.Blob {
	return s.s.Get(ctx, id.Hex())
}

func (s *Storage) Add(ctx context.Context, b blob.Blob) (*Node, error) {

	// TODO: avoid full memory copy

	data, err := blob.Bytes(b)
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}
	id := NewIDFromBytes(data)
	n, err := s.Node(ctx, id)
	if err == nil {
		return n, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	key := id.Hex()
	if err := s.s.Set(ctx, key, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	attrs, err := file.ResolveAttrs(blob.FromBytes(data), s.fileOpts...)
	if err != nil {
		return nil, err
	}
	return &Node{
		ID:         id,
		Attributes: *attrs,
	}, s.ds.Set(ctx, key, *attrs)
}

// TODO: remover

func (s *Storage) Iter(ctx context.Context, after ID) rx.Iter[rx.Pair[ID, file.Attributes]] {
	return rx.Map(s.ds.Iter(ctx, after.Hex()), func(
		p rx.Pair[string, file.Attributes]) rx.Pair[ID, file.Attributes] {

		id, _ := ParseID(p.Key)
		return rx.Pair[ID, file.Attributes]{Key: id, Value: p.Value}
	})
}
