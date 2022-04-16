package cas

import (
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
)

type Node struct {
	ID ID `json:"id"`
	file.Attributes

	s *Storage
}

func (n Node) File() blob.Blob {
	return n.s.s.Get(context.Background(), n.ID.Hex())
}
