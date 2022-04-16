package cas

import (
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
)

type Node struct {
	ID ID `json:"id"`
	file.Attributes

	s blob.Storage
}

func (n Node) File() blob.Blob {
	return n.s.Get(context.Background(), n.ID.Hex())
}
