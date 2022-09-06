package redis

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/koud-fi/pkg/blob"

	"github.com/go-redis/redis/v8"
)

type Options = redis.Options

func Open(opts *Options) *redis.Client {
	return redis.NewClient(opts)
}

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	rc         *redis.Client
	keyPrefix  string
	expiration time.Duration
}

type Option func(s *Storage)

func KeyPrefix(prefix string) Option      { return func(s *Storage) { s.keyPrefix = prefix } }
func Expiration(dur time.Duration) Option { return func(s *Storage) { s.expiration = dur } }

func NewStorage(rc *redis.Client, opt ...Option) *Storage {
	s := Storage{rc: rc}
	for _, opt := range opt {
		opt(&s)
	}
	return &s
}

func (s *Storage) Get(ctx context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		data, err := s.rc.Get(ctx, s.key(ref)).Bytes()
		if err != nil {
			if err == redis.Nil {
				return nil, os.ErrNotExist
			}
			return nil, err
		}
		return data, nil
	})
}

func (s *Storage) Set(ctx context.Context, ref string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.rc.Set(ctx, s.key(ref), data, s.expiration).Err()
}

func (s *Storage) Delete(ctx context.Context, refs ...string) error {

	// ???

	panic("TODO")
}

func (s *Storage) key(ref string) string { return s.keyPrefix + ref }
