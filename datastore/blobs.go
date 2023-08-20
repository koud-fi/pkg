package datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

// TODO: make the (un)marshaler configurable

type blobsKV[T any] struct {
	blobs blob.Storage
}

func BlobsKV[T any](bs blob.Storage) KV[T] {
	return &blobsKV[T]{blobs: bs}
}

func (bkv *blobsKV[T]) Get(ctx context.Context, key string) (v T, err error) {
	err = blob.Unmarshal(json.Unmarshal, bkv.blobs.Get(ctx, key), &v)
	return
}

func (bkv *blobsKV[T]) Put(ctx context.Context, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if currData, err := blob.Bytes(bkv.blobs.Get(ctx, key)); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if bytes.Equal(data, currData) {
		return nil
	}
	return bkv.blobs.Set(ctx, key, bytes.NewReader(data))
}

func (bkv *blobsKV[T]) Delete(ctx context.Context, keys ...string) error {
	return bkv.blobs.Delete(ctx, keys...)
}

type sortedBlobsKV[T any] struct {
	blobsKV[T]
	sorted blob.SortedStorage
}

func SortedBlobsKV[T any](
	sbs blob.SortedStorage, keyFn func(T) (string, error),
) SortedKV[T] {
	return &sortedBlobsKV[T]{
		blobsKV: blobsKV[T]{sbs},
		sorted:  sbs,
	}
}

func (sbkv sortedBlobsKV[T]) Iter(ctx context.Context, after string) rx.Iter[rx.Pair[string, T]] {

	// TODO: lazy iterator creation

	it := sbkv.sorted.Iter(ctx, after)
	return rx.MapErr(it, func(rb blob.RefBlob) (rx.Pair[string, T], error) {
		var v T
		err := blob.Unmarshal(json.Unmarshal, rb.Blob, &v)
		return rx.NewPair(rb.Ref, v), err
	})
}

type blobsTable[T any] struct {
	kv    KV[T]
	keyFn func(T) (string, error)
}

func BlobsTable[T any](bs blob.Storage, keyFn func(T) (string, error)) Table[T] {
	return &blobsTable[T]{BlobsKV[T](bs), keyFn}
}

func (bt *blobsTable[T]) Get(ctx context.Context) func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
	return func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
		ref, err := bt.keyFn(key)
		if err != nil {
			return rx.Pair[T, rx.Maybe[T]]{}, err
		}
		v, err := bt.kv.Get(ctx, ref)
		if err != nil {
			if os.IsNotExist(err) {
				return rx.NewPair(key, rx.None[T]()), nil
			}
			return rx.Pair[T, rx.Maybe[T]]{}, err
		}
		return rx.NewPair(key, rx.Just(v)), nil
	}
}

func (bt *blobsTable[T]) Put(ctx context.Context) func(value T) (T, error) {
	return func(value T) (v T, _ error) {
		ref, err := bt.keyFn(v)
		if err != nil {
			return v, err
		}
		return value, bt.kv.Put(ctx, ref, value)
	}
}

func (bt *blobsTable[T]) Delete(ctx context.Context) func(key T) error {
	return func(key T) error {
		ref, err := bt.keyFn(key)
		if err != nil {
			return err
		}
		return bt.kv.Delete(ctx, ref)
	}
}

// TODO: sorted blobs table implementation
