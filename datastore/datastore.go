package datastore

import (
	"bytes"
	"context"

	"github.com/koud-fi/pkg/blob"
)

type Store[T any] struct {
	s blob.Storage
	c Codec
}

type Codec struct {
	Marshal   blob.MarshalFunc
	Unmarshal blob.UnmarshalFunc
}

func New[T any](s blob.Storage, c Codec) *Store[T] {
	return &Store[T]{s: s}
}

func (s *Store[T]) Get(ctx context.Context, key string) (T, error) {
	var v T
	err := blob.Unmarshal(s.c.Unmarshal, s.s.Get(ctx, key), &v)
	return v, err
}

func (s *Store[T]) Set(ctx context.Context, key string, v T) error {
	b, err := s.c.Marshal(v)
	if err != nil {
		return err
	}
	return s.s.Set(ctx, key, bytes.NewReader(b))
}

func (s *Store[T]) Update(ctx context.Context, key string, fn func(v T) (T, error)) error {
	b1, err := blob.Bytes(s.s.Get(ctx, key))
	if err != nil {
		return err
	}
	var v T
	if err := s.c.Unmarshal(b1, &v); err != nil {
		return err
	}
	if v, err = fn(v); err != nil {
		return err
	}
	b2, err := s.c.Marshal(v)
	if err != nil {
		return err
	}
	if bytes.Compare(b1, b2) == 0 {
		return nil
	}
	return s.s.Set(ctx, key, bytes.NewReader(b2))
}

func (s *Store[T]) Delete(ctx context.Context, key ...string) error {
	return s.s.Delete(ctx, key...)
}

// TODO: enumeration