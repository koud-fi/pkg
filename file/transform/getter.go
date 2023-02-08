package transform

import (
	"context"
	"io"

	"github.com/koud-fi/pkg/blob"
)

const (
	Image blob.Domain = "image"
	//Video blob.Domain = "video"
)

func ImageGetter(g blob.Getter) blob.Getter {
	return blob.GetterFunc(func(ctx context.Context, ref blob.Ref) blob.Blob {
		return blob.Func(func() (io.ReadCloser, error) {
			ref = ref.Ref()
			p, err := ParseParams(string(ref.Domain()))
			if err != nil {
				return nil, err
			}
			return ToImage(g.Get(ctx, ref.Ref()), p).Open()
		})
	})
}
