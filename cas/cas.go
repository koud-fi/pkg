package cas

import (
	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/grf"
)

type Storage struct {
	s blob.Storage
	g *grf.Graph
}

func New(s blob.Storage, g *grf.Graph) *Storage {
	return &Storage{s: s, g: g}
}

// TODO
