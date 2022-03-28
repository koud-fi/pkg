package blobcache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/cache"
)

//lint:ignore ST1012 if returned from resolve; file is not written to storage.
var NoCache = errors.New("no cache")

type Cache struct {
	*cache.Cache
	backend
}

func New(s blob.Storage) *Cache {
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
		if err := c.Cache.Resolve(key, func() (int64, error) {
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
	s blob.Storage
}

func (b *backend) Has(key string) (bool, error) {

	// TODO: use blob.Statter

	panic("TODO")
}

func (b *backend) Delete(key string) error {

	// ???

	panic("TODO")
}

func (b *backend) Keys(fn func(key string, size int64)) error {

	// ???

	panic("TODO")
}
