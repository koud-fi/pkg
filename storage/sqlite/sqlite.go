package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0700)); err != nil {
		return nil, err
	}
	return sql.Open("sqlite3", fmt.Sprintf("file:%s", path))
}

type Storage struct {
	db *sql.DB
}

var _ blob.Storage = (*Storage)(nil)

func NewStorage(db *sql.DB, table string) *Storage {

	// TODO: init table

	return &Storage{db: db}
}

func (s *Storage) Get(ctx context.Context, ref string) blob.Blob {

	// ???

	panic("TODO")
}

func (s *Storage) Set(ctx context.Context, ref string, r io.Reader) error {

	// ???

	panic("TODO")
}

func (s *Storage) Iter(ctx context.Context, after string) rx.Iter[blob.RefBlob] {

	// ???

	panic("TODO")
}

func (s *Storage) Delete(ctx context.Context, refs ...string) error {

	// ???

	panic("TODO")
}
