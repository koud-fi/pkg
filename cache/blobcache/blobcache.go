package blobcache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/cache"
	"github.com/koud-fi/pkg/rx"
)

//lint:ignore ST1012 if returned from resolve; data is not written to storage.
var NoCache = errors.New("no cache")

type Cache struct {
	*cache.Cache
	backend
}

func New(s blob.SortedStorage) *Cache {
	b := backend{s: s}
	return &Cache{cache.New(&b), b}
}

func (c *Cache) Resolve(ctx context.Context, key string, b blob.Blob) blob.Blob {
	return blob.Func(func() (io.ReadCloser, error) {
		var (
			digest = sha256.Sum256([]byte(key))
			key    = hex.EncodeToString(digest[:])
			out    io.ReadCloser
		)
		if err := c.Cache.Resolve(ctx, key, func() (int64, error) {
			rc, err := b.Open()
			if err != nil {
				if err == NoCache {
					out = rc
					return 0, nil
				}
				return 0, err
			}
			defer rc.Close()

			// TODO: resolve size correctly

			return 0, c.s.Set(ctx, key, rc)

		}); err != nil {
			return nil, err
		}
		if out != nil {
			return out, nil
		}
		return c.s.Get(ctx, key).Open()
	})
}

type backend struct {
	s blob.SortedStorage
}

func (b *backend) Has(ctx context.Context, key string) (bool, error) {
	if err := blob.Error(b.s.Get(ctx, key)); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

func (b *backend) Delete(ctx context.Context, key string) error {
	return b.s.Delete(ctx, key)
}

func (b *backend) Keys(ctx context.Context) rx.Iter[rx.Pair[string, int64]] {
	return rx.Map(b.s.Iter(ctx, ""), func(br blob.RefBlob) rx.Pair[string, int64] {
		return rx.Pair[string, int64]{Key: br.Ref, Value: -1} // TODO: resolve blob sizes
	})
}
