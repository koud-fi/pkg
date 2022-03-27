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

func Fetch(ctx context.Context, ref string) blob.Blob {
	panic("TODO")
}

func Receive(ctx context.Context, ref string, r io.Reader) error {
	panic("TODO")
}

func Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {
	panic("TODO")
}

func Stat(ctx context.Context, refs []string, fn func(string, int64) error) error {
	panic("TODO")
}

func Remove(ctx context.Context, refs ...string) error {
	panic("TODO")
}
