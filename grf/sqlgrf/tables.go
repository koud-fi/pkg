package sqlgrf

import (
	"fmt"

	"github.com/koud-fi/pkg/grf"
)

const (
	nodeTableSuffix      = "_node"
	edgeTypeTableSuffix  = "_edgetype"
	edgeTableSuffix      = "_edge"
	edgeCountTableSuffix = "_edgecount"
	keymapTableSuffix    = "_keymap"
)

type tables struct {
	nodes      string
	edgeTypes  string
	edges      string
	edgeCounts string
}

func (s *store) tables(nt grf.NodeType) (tables, error) {
	s.tablesMu.Lock()
	defer s.tablesMu.Unlock()

	if t, ok := s.tablesMap[nt]; ok {
		return t, nil
	}
	t := tables{
		nodes:      string(nt) + nodeTableSuffix,
		edgeTypes:  string(nt) + edgeTypeTableSuffix,
		edges:      string(nt) + edgeTableSuffix,
		edgeCounts: string(nt) + edgeCountTableSuffix,
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id   INTEGER  PRIMARY KEY AUTOINCREMENT,
			data BLOB     NULL,
			ts   DATETIME NOT NULL
		)
	`, t.nodes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.nodes, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			type  TEXT    NOT NULL UNIQUE
		)
	`, t.edgeTypes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.edgeTypes, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			from_id  INTEGER NOT NULL,
			type_id  INTEGER NOT NULL,
			to_id    INTEGER NOT NULL,
			sequence INTEGER NOT NULL,
			data     BLOB NULL,

			PRIMARY KEY(from_id, type_id, to_id),
			FOREIGN KEY(from_id) REFERENCES %s(id) ON DELETE CASCADE,
			FOREIGN KEY(type_id) REFERENCES %s(id) ON DELETE RESTRICT
		)
	`, t.edges, t.nodes, t.edgeTypes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_sequence ON %s (sequence)
	`, t.edges, t.edges)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s sequence index: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			from_id INTEGER NOT NULL,
			type_id INTEGER NOT NULL,
			count INTEGER NOT NULL DEFAULT 0,

			PRIMARY KEY (from_id, type_id)
			FOREIGN KEY(from_id) REFERENCES %s(id) ON DELETE CASCADE,
			FOREIGN KEY(type_id) REFERENCES %s(id) ON DELETE RESTRICT
		)
	`, t.edgeCounts, t.nodes, t.edgeTypes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.edgeCounts, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS %s_insert AFTER INSERT ON %s 
		BEGIN
			INSERT INTO %s (from_id, type_id, count) VALUES (NEW.from_id, NEW.type_id, 1)
			ON CONFLICT(from_id, type_id) DO
				UPDATE SET count = count + 1;
		END
	`, t.edges, t.edges, t.edgeCounts)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s insert trigger: %w", t.edges, err)
	}

	/*
		INSERT INTO %s (from_id, type_id, to_id, sequence, data) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(from_id, type_id, to_id) DO
			UPDATE SET sequence = excluded.sequence, data = excluded.data
	*/

	// TODO: trigger for edge deletion counting

	return t, nil
}

func (m *mapper) table(nt grf.NodeType) (string, error) {
	m.tableMu.Lock()
	defer m.tableMu.Unlock()

	if t, ok := m.tableMap[nt]; ok {
		return t, nil
	}
	t := string(nt) + keymapTableSuffix
	if _, err := m.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			key TEXT PRIMARY KEY,
			id  INT NOT NULL
		)
	`, t)); err != nil {
		return "", err
	}
	return t, nil
}
