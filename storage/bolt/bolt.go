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

var (
	dataSuffix = []byte(":data")
	statSuffix = []byte(":stat")
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
	bucket []byte
	db     *DB
}

func NewStorage(db *DB, bucket string) *Storage {
	return &Storage{[]byte(bucket), db}
}

func (s Storage) Get(_ context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() (data []byte, err error) {
		err = s.db.View(func(tx *bolt.Tx) error {
			return wrapTx(tx).useBucket(s.bucket, true, false, func(b *bolt.Bucket) error {
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

func (s Storage) Set(_ context.Context, ref string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		key := []byte(ref)
		return wrapTx(tx).
			useBucket(s.bucket, true, false, func(b *bolt.Bucket) error {
				return b.Put(key, data)
			}).
			useBucket(s.bucket, false, false, func(b *bolt.Bucket) error {
				return b.Put(key, []byte(strconv.Itoa(len(data))))
			}).err
	})
}

func (s Storage) Delete(_ context.Context, refs ...string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, ref := range refs {
			key := []byte(ref)
			return wrapTx(tx).
				useBucket(s.bucket, false, true, func(b *bolt.Bucket) error {
					return b.Delete(key)
				}).
				useBucket(s.bucket, true, true, func(b *bolt.Bucket) error {
					return b.Delete(key)
				}).
				err
		}
		return nil
	})
}

func (s *Storage) Iter(ctx context.Context, state rx.Lens[string]) rx.Iter[blob.RefBlob] {
	return &iter{
		s:      s,
		ctx:    ctx,
		bucket: append(s.bucket, statSuffix...),
		state:  state,
	}
}

type iter struct {
	s      *Storage
	ctx    context.Context
	bucket []byte
	state  rx.Lens[string]

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
		if b := it.tx.Bucket(it.bucket); b != nil {
			after, err := it.state.Get()
			if err != nil {
				it.err = err
				return false
			}
			it.c = b.Cursor()
			if len(after) == 0 {
				it.k, _ = it.c.First()
			} else {
				it.c.Seek([]byte(after))
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
			if it.err = it.state.Set(string(it.k)); it.err != nil {
				return false
			}
			it.k, _ = it.c.Next()
		}
		return it.k != nil
	}
}

func (it iter) Value() blob.RefBlob {
	ref := string(it.k)
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

type txWrapper struct {
	tx  *bolt.Tx
	err error
}

func wrapTx(tx *bolt.Tx) *txWrapper { return &txWrapper{tx: tx} }

func (w *txWrapper) useBucket(
	bucket []byte, isData bool, skipNil bool, fn func(b *bolt.Bucket) error,
) *txWrapper {
	var name []byte
	if isData {
		name = append(bucket, dataSuffix...)
	} else {
		name = append(bucket, statSuffix...)
	}
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
