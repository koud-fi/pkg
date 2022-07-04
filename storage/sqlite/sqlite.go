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
	db    *sql.DB
	table string
}

var _ blob.Storage = (*Storage)(nil)

func NewStorage(db *sql.DB, table string) *Storage {
	if _, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			ref       TEXT
			data_size INTEGER
			data      BLOB
		)
	`, table)); err != nil {
		panic("failed to create sqlite table: " + err.Error())
	}
	return &Storage{db: db, table: table}
}

func (s *Storage) Get(ctx context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		var buf []byte
		return buf, s.db.QueryRowContext(ctx, fmt.Sprintf(`
			SELECT data FROM %s
			WHERE ref = ? 
		`, s.table), ref).Scan()
	})
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
