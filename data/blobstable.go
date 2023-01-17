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
	refFn func(T) (string, error)
}

func BlobsTable[T any](bs blob.Storage, refFn func(T) (string, error)) Table[T] {
	return &blobsTable[T]{blobs: bs, refFn: refFn}
}

func (bt *blobsTable[T]) Get(ctx context.Context, keys rx.Iter[T]) rx.Iter[rx.Pair[T, rx.Maybe[T]]] {
	refKeys := rx.PluckErr(keys, bt.refFn)
	return rx.MapErr(refKeys, func(p rx.Pair[string, T]) (rx.Pair[T, rx.Maybe[T]], error) {
		var v T
		if err := blob.Unmarshal(json.Unmarshal, bt.blobs.Get(ctx, p.Key()), &v); err != nil {
			if os.IsNotExist(err) {
				return rx.NewPair(p.Value(), rx.None[T]()), nil
			}
			return rx.Pair[T, rx.Maybe[T]]{}, err
		}
		return rx.NewPair(p.Value(), rx.Just(v)), nil
	})
}

func (bt *blobsTable[T]) Put(ctx context.Context, values rx.Iter[T]) rx.Iter[T] {
	return rx.MapErr(rx.PluckErr(values, bt.refFn), func(p rx.Pair[string, T]) (v T, _ error) {
		data, err := json.Marshal(p.Value())
		if err != nil {
			return v, err
		}
		return p.Value(), bt.blobs.Set(ctx, p.Key(), bytes.NewReader(data))
	})
}

func (bt *blobsTable[T]) Delete(ctx context.Context, keys rx.Iter[T]) error {
	return rx.UseSlice(rx.MapErr(keys, bt.refFn), func(refs []string) error {
		return bt.blobs.Delete(ctx, refs...)
	})
}
