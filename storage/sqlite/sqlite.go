package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: blob.Storage is a completely wrong level of abstraction for this, implement this as a data.Table

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0700)); err != nil {
		return nil, err
	}
	return sql.Open("sqlite3", fmt.Sprintf("file:%s", path))
}

/*
type Storage struct {
	db    *sql.DB
	table string
}

var _ blob.SortedStorage = (*Storage)(nil)

func NewStorage(db *sql.DB, table string) *Storage {
	if _, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			ref       TEXT PRIMARY KEY,
			data_size INTEGER,
			data      BLOB
		)
	`, table)); err != nil {
		panic("failed to create sqlite table: " + err.Error())
	}
	return &Storage{db: db, table: table}
}

func (s *Storage) Get(ctx context.Context, ref blob.Ref) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		var buf []byte
		if err := s.db.QueryRowContext(ctx, fmt.Sprintf(`
			SELECT data FROM %s
			WHERE ref = ?
		`, s.table), ref).Scan(&buf); err != nil {
			if err == sql.ErrNoRows {
				return nil, os.ErrNotExist
			}
			return nil, err
		}
		return buf, nil
	})
}

func (s *Storage) Set(ctx context.Context, ref blob.Ref, r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, fmt.Sprintf(`
		INSERT INTO %s (ref, data_size, data) VALUES (?, ?, ?)
		ON CONFLICT(ref) DO UPDATE SET
			data_size = excluded.data_size,
			data = excluded.data
	`, s.table), ref, len(buf), buf)
	return err
}

func (s *Storage) Iter(ctx context.Context, after blob.Ref) rx.Iter[blob.RefBlob] {
	return &iter{s: s, ctx: ctx, after: after.String()}
}

type iter struct {
	s     *Storage
	ctx   context.Context
	after string

	rows *sql.Rows
	ref  string
	data []byte
	err  error
}

func (it *iter) Next() bool {
	if it.rows == nil {
		if it.rows, it.err = it.s.db.QueryContext(it.ctx, fmt.Sprintf(`
			SELECT ref, data FROM %s WHERE ref > ?
			ORDER BY ref ASC
		`, it.s.table), it.after); it.err != nil {
			return false
		}
	}
	if !it.rows.Next() {
		return false
	}
	it.err = it.rows.Scan(&it.ref, &it.data)
	return it.err == nil
}

func (it iter) Value() blob.RefBlob {
	return blob.RefBlob{
		Ref:  blob.ParseRef(it.ref),
		Blob: blob.FromBytes(it.data),
	}
}

func (it iter) Close() error {
	var closeErr error
	if it.rows != nil {
		closeErr = it.rows.Close()
	}
	if it.err != nil {
		return closeErr
	}
	return it.err
}

func (s *Storage) Delete(ctx context.Context, refs ...blob.Ref) error {

	// TODO: optimize (use a single query)

	for _, ref := range refs {
		if _, err := s.db.ExecContext(ctx, fmt.Sprintf(`
			DELETE FROM %s WHERE ref = ?
		`, s.table), ref); err != nil {
			return err
		}
	}
	return nil
}
*/
