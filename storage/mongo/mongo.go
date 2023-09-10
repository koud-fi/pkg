package mongo

import "context"

type KV[T any] struct {
	// TODO
}

func (kv *KV[T]) Get(ctx context.Context, key string) (T, error) {
	panic("TODO")
}

func (kv *KV[T]) Put(ctx context.Context, key string, value T) error {
	panic("TODO")
}

func (kv *KV[T]) Delete(ctx context.Context, keys ...string) error {
	panic("TODO")
}
