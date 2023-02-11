package bolt

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"

	bolt "go.etcd.io/bbolt"
)

type DB = bolt.DB

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0700)); err != nil {
		return nil, err
	}
	return bolt.Open(path, os.FileMode(0600), nil)
}

var _ blob.SortedStorage = (*Storage)(nil)

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

func (s Storage) Get(_ context.Context, ref blob.Ref) blob.Blob {
	return blob.ByteFunc(func() (data []byte, err error) {
		err = s.db.View(func(tx *bolt.Tx) error {
			return wrapTx(tx).useBucket(s.dataBucket, false, func(b *bolt.Bucket) error {
				if b == nil {
					return os.ErrNotExist
				}
				buf := b.Get(ref.Bytes())
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

func (s Storage) Set(_ context.Context, ref blob.Ref, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		key := ref.Bytes()
		return wrapTx(tx).
			useBucket(s.dataBucket, false, func(b *bolt.Bucket) error {
				return b.Put(key, data)
			}).
			useBucket(s.statBucket, false, func(b *bolt.Bucket) error {
				return b.Put(key, []byte(strconv.Itoa(len(data))))
			}).err
	})
}

func (s *Storage) Iter(ctx context.Context, after blob.Ref) rx.Iter[blob.RefBlob] {
	return &iter{s: s, ctx: ctx, after: after.Bytes()}
}

type iter struct {
	s     *Storage
	ctx   context.Context
	after []byte

	tx  *bolt.Tx
	c   *bolt.Cursor
	k   []byte
	err error
}

func (it *iter) Next() bool {
	var init bool
	if it.c == nil {
		it.tx, it.err = it.s.db.Begin(false)
		if it.err != nil {
			return false
		}
		if b := it.tx.Bucket(it.s.statBucket); b != nil {
			it.c = b.Cursor()

			if len(it.after) == 0 {
				it.k, _ = it.c.First()
			} else {
				it.c.Seek(it.after)
			}
		}
		init = true
	}
	if it.k == nil {
		return false
	}
	select {
	case <-it.ctx.Done():
		it.err = it.ctx.Err()
		return false
	default:
		if !init {
			it.k, _ = it.c.Next()
		}
		return it.k != nil
	}
}

func (it iter) Value() blob.RefBlob {
	ref := blob.ParseRef(string(it.k))
	return blob.RefBlob{
		Ref:  ref,
		Blob: it.s.Get(it.ctx, ref),
	}
}

func (it iter) Close() error {
	var rollbackErr error
	if it.tx != nil {
		rollbackErr = it.tx.Rollback()
	}
	if it.err != nil {
		return rollbackErr
	}
	return it.err
}

func (s Storage) Delete(_ context.Context, refs ...blob.Ref) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, ref := range refs {
			key := ref.Bytes()
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
