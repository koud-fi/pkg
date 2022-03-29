package cache

import (
	"sync/atomic"

	"golang.org/x/sync/singleflight"
)

type Backend interface {
	Has(key string) (bool, error)
	Delete(key string) error
	Keys(fn func(key string, size int64)) error
}

type Cache struct {
	b Backend
	g singleflight.Group

	// TODO: capacity limiting
	// TODO: expire times
	// TODO: key cache

	count     int64
	totalSize int64
}

func New(b Backend) *Cache { return &Cache{b: b} }

func (c *Cache) Resolve(key string, fn func() (int64, error)) error {
	_, err, _ := c.g.Do(key, func() (any, error) {
		if ok, err := c.b.Has(key); err != nil {
			return nil, err
		} else if ok {
			return nil, nil
		}
		size, err := fn()
		if err != nil {
			return nil, err
		}
		atomic.AddInt64(&c.count, 1)
		atomic.AddInt64(&c.totalSize, size)
		return nil, nil
	})
	return err
}
