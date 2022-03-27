package bolt

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/koud-fi/pkg/blob"

	bolt "go.etcd.io/bbolt"
)

type DB = bolt.DB

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0700)); err != nil {
		return nil, err
	}
	return bolt.Open(path, os.FileMode(0600), nil)
}

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	db         *DB
	dataBucket []byte
	statBucket []byte
}

func NewStorage(db *DB, bucket string) *Storage {
	return &Storage{db: db,
		dataBucket: []byte(bucket + ":data"),
		statBucket: []byte(bucket + ":stat"),
	}
}

func (s Storage) Fetch(ctx context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() (data []byte, err error) {
		err = s.db.View(func(tx *bolt.Tx) error {
			return wrapTx(tx).useBucket(s.dataBucket, false, func(b *bolt.Bucket) error {
				if b == nil {
					return os.ErrNotExist
				}
				buf := b.Get([]byte(ref))
				if buf == nil {
					return os.ErrNotExist
				}
				data = make([]byte, len(buf))
				copy(data, buf)
				return nil
			}).err
		})
		return
	})
}

func (s Storage) Receive(ctx context.Context, ref string, r io.Reader) error {

	// ???

	panic("TODO")
}

func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {

	// ???

	panic("TODO")
}

func (s Storage) Stat(ctx context.Context, refs []string, fn func(string, int64) error) error {

	// ???

	panic("TODO")
}

func (s Storage) Remove(ctx context.Context, refs ...string) error {

	// ???

	panic("TODO")
}

type txWrapper struct {
	tx  *bolt.Tx
	err error
}

func wrapTx(tx *bolt.Tx) *txWrapper { return &txWrapper{tx: tx} }

func (w *txWrapper) useBucket(name []byte, skipNil bool, fn func(b *bolt.Bucket) error) *txWrapper {
	if w.err != nil {
		return w
	}
	var b *bolt.Bucket
	if w.tx.Writable() {
		if b, w.err = w.tx.CreateBucketIfNotExists(name); w.err != nil {
			return w
		}
	} else {
		b = w.tx.Bucket(name)
	}
	if !skipNil || b != nil {
		w.err = fn(b)
	}
	return w
}
