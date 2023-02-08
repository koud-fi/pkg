package fetch

import (
	"context"

	"github.com/koud-fi/pkg/blob"
)

const (
	Http  blob.Domain = "http"
	Https blob.Domain = "https"
)

func Getter(reqFn func(context.Context, blob.Ref) blob.Blob) blob.Getter {
	return blob.GetterFunc(func(ctx context.Context, ref blob.Ref) blob.Blob {
		return reqFn(ctx, ref)
	})
}
