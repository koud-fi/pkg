package data

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

type blobsTable[T any] struct {
	blobs blob.Storage
	refFn func(T) (blob.Ref, error)
}

func BlobsTable[T any](bs blob.Storage, refFn func(T) (blob.Ref, error)) Table[T] {
	return &blobsTable[T]{blobs: bs, refFn: refFn}
}

func (bt *blobsTable[T]) Get(ctx context.Context) func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
	return func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
		ref, err := bt.refFn(key)
		if err != nil {
			return rx.Pair[T, rx.Maybe[T]]{}, err
		}
		var v T
		if err := blob.Unmarshal(json.Unmarshal, bt.blobs.Get(ctx, ref), &v); err != nil {
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
		ref, err := bt.refFn(value)
		if err != nil {
			return v, err
		}
		data, err := json.Marshal(value)
		if err != nil {
			return v, err
		}
		if currData, err := blob.Bytes(bt.blobs.Get(ctx, ref)); err != nil {
			if !os.IsNotExist(err) {
				return v, err
			}
		} else if bytes.Equal(data, currData) {
			return value, nil
		}
		return value, bt.blobs.Set(ctx, ref, bytes.NewReader(data))
	}
}

func (bt *blobsTable[T]) Delete(ctx context.Context) func(key T) error {
	return func(key T) error {
		ref, err := bt.refFn(key)
		if err != nil {
			return err
		}
		return bt.blobs.Delete(ctx, ref)
	}
}
