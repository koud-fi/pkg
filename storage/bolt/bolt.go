package bolt

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"

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

func (s Storage) Fetch(_ context.Context, ref string) blob.Blob {
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

func (s Storage) Receive(_ context.Context, ref string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		key := []byte(ref)
		return wrapTx(tx).
			useBucket(s.dataBucket, false, func(b *bolt.Bucket) error {
				return b.Put(key, data)
			}).
			useBucket(s.statBucket, false, func(b *bolt.Bucket) error {
				return b.Put(key, []byte(strconv.Itoa(len(data))))
			}).err
	})
}

func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {
	tx, err := s.db.Begin(false)
	if err != nil {
		return err
	}
	b := tx.Bucket(s.statBucket)
	if b == nil {
		return nil
	}
	var (
		c          *bolt.Cursor
		afterBytes = []byte(after)
		k, v       []byte
	)
	if after == "" {
		k, v = c.First()
	} else {
		c.Seek(afterBytes)
		k, v = c.Next()
	}
	for k != nil {
		select {
		case <-ctx.Done():
			return nil
		default:
			size, err := strconv.ParseInt(string(v), 10, 64)
			if err != nil {
				return err
			}
			fn(string(k), size)
			k, v = c.Next()
		}
	}
	return nil
}

func (s Storage) Stat(_ context.Context, refs []string, fn func(string, int64) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		for _, ref := range refs {
			if err := wrapTx(tx).useBucket(s.statBucket, true, func(b *bolt.Bucket) error {
				buf := b.Get([]byte(ref))
				if buf == nil {
					return nil
				}
				size, err := strconv.ParseInt(string(buf), 10, 64)
				if err != nil {
					return err
				}
				return fn(ref, size)
			}).err; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s Storage) Remove(_ context.Context, refs ...string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, ref := range refs {
			key := []byte(ref)
			return wrapTx(tx).
				useBucket(s.statBucket, true, func(b *bolt.Bucket) error {
					return b.Delete(key)
				}).
				useBucket(s.dataBucket, true, func(b *bolt.Bucket) error {
					return b.Delete(key)
				}).
				err
		}
		return nil
	})
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
