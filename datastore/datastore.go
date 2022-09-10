package datastore

import (
	"bytes"
	"context"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

type Store[T any] struct {
	s blob.SortedStorage
	c Codec[T]
}

func New[T any](s blob.SortedStorage, c Codec[T]) *Store[T] {
	return &Store[T]{s: s, c: c}
}

func (s *Store[T]) Get(ctx context.Context, key string) (T, error) {
	return s.c.Unmarshal(s.s.Get(ctx, key))
}

func (s *Store[T]) Put(ctx context.Context, key string, v T) error {
	rc, err := s.c.Marshal(v).Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return s.s.Set(ctx, key, rc)
}

// Set is an alias to Put.
func (s *Store[T]) Set(ctx context.Context, key string, v T) error { return s.Put(ctx, key, v) }

func (s *Store[T]) Update(ctx context.Context, key string, fn func(v T) (T, error)) error {
	b1, err := blob.Bytes(s.s.Get(ctx, key))
	if err != nil {
		return err
	}
	v, err := s.c.Unmarshal(blob.FromBytes(b1))
	if err != nil {
		return err
	}
	if v, err = fn(v); err != nil {
		return err
	}
	b2, err := blob.Bytes(s.c.Marshal(v))
	if err != nil {
		return err
	}
	if bytes.Equal(b1, b2) {
		return nil
	}
	return s.s.Set(ctx, key, bytes.NewReader(b2))
}

func (s *Store[T]) Iter(ctx context.Context, after string) rx.Iter[rx.Pair[string, T]] {
	return rx.MapErr(s.s.Iter(ctx, after), func(b blob.RefBlob) (rx.Pair[string, T], error) {
		v, err := s.c.Unmarshal(b)
		return rx.Pair[string, T]{Key: b.Ref, Value: v}, err
	})
}

func (s *Store[T]) Delete(ctx context.Context, key ...string) error {
	return s.s.Delete(ctx, key...)
}
