package cas

import (
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
)

type Node struct {
	ID ID `json:"id"`
	NodeData

	s blob.Storage
}

type NodeData struct {
	file.Attributes
}

func (n Node) File() blob.Blob {
	return n.s.Get(context.Background(), n.ID.Hex())
}
