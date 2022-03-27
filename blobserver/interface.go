package blobserver

import (
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/pk"
)

type Fetcher interface {
	Fetch(context.Context, pk.Ref) blob.Blob
}

type FetchFunc func(context.Context, pk.Ref) blob.Blob

func (fn FetchFunc) Fetch(ctx context.Context, ref pk.Ref) blob.Blob {
	return fn(ctx, ref)
}
