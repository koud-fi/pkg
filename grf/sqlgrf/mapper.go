package sqlgrf

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/koud-fi/pkg/grf"
)

type mapper struct {
	db   *sql.DB
	init sync.Once

	tableMu  sync.Mutex
	tableMap map[grf.NodeType]string
}

func NewMapper(db *sql.DB) grf.Mapper {
	return &mapper{db: db}
}

func (m *mapper) Map(nt grf.NodeType, key string) (grf.ID, error) {
	t, err := m.table(nt)
	if err != nil {
		return 0, err
	}
	var id grf.ID
	if err := m.db.QueryRow(fmt.Sprintf(`
		SELECT id FROM %s WHERE key = ?
	`, t), key).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, grf.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (m *mapper) SetMapping(nt grf.NodeType, key string, id grf.ID) error {
	t, err := m.table(nt)
	if err != nil {
		return err
	}
	_, err = m.db.Exec(fmt.Sprintf(`
		INSERT INTO %s (key, id) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET id = excluded.id
	`, t), key, id)
	return err
}

func (m *mapper) DeleteMapping(nt grf.NodeType, key ...string) error {
	t, err := m.table(nt)
	if err != nil {
		return err
	}
	if len(key) == 0 {
		return nil
	}
	_, err = m.db.Exec(fmt.Sprintf(`
		DELETE FROM %s WHERE key IN (%s)
	`, t, strings.Join(key, ",")))
	return err
}
