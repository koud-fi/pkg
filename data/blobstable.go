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
	keyFn func(T) (string, error)
}

func BlobsTable[T any](bs blob.Storage, keyFn func(T) (string, error)) Table[T] {
	return &blobsTable[T]{bs, keyFn}
}

func (bt *blobsTable[T]) Get(ctx context.Context) func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
	return func(key T) (rx.Pair[T, rx.Maybe[T]], error) {
		ref, err := bt.keyFn(key)
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
		ref, err := bt.keyFn(v)
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
		ref, err := bt.keyFn(key)
		if err != nil {
			return err
		}
		return bt.blobs.Delete(ctx, ref)
	}
}

type sortedBlobsTable[T any] struct {
	blobsTable[T]
	sorted blob.SortedStorage
}

func SortedBlobsTable[T any](
	sbs blob.SortedStorage, keyFn func(T) (string, error),
) SortedTable[T] {
	return &sortedBlobsTable[T]{
		blobsTable: blobsTable[T]{sbs, keyFn},
		sorted:     sbs,
	}
}

func (sbt sortedBlobsTable[T]) Iter(ctx context.Context, after T) rx.Iter[T] {

	// TODO: lazy iterator creation

	afterRef, err := sbt.keyFn(after)
	if err != nil {
		return rx.Error[T](err)
	}
	it := sbt.sorted.Iter(ctx, afterRef)
	return rx.MapErr(it, func(rb blob.RefBlob) (v T, err error) {
		err = blob.Unmarshal(json.Unmarshal, rb.Blob, &v)
		return
	})
}
