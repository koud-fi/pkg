package sqlgrf

import (
	"fmt"

	"github.com/koud-fi/pkg/grf"
)

const (
	nodeTableSuffix     = "_node"
	edgeTableSuffix     = "_edge"
	edgeInfoTableSuffix = "_edgeinfo"
	keymapTableSuffix   = "_keymap"
)

type tables struct {
	nodes     string
	edges     string
	edgeInfos string
}

func (s *store) tables(nt grf.NodeType) (tables, error) {
	s.tablesMu.Lock()
	defer s.tablesMu.Unlock()

	if t, ok := s.tablesMap[nt]; ok {
		return t, nil
	}
	t := tables{
		nodes:     string(nt) + nodeTableSuffix,
		edges:     string(nt) + edgeTableSuffix,
		edgeInfos: string(nt) + edgeInfoTableSuffix,
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id      INTEGER  PRIMARY KEY AUTOINCREMENT,
			data    BLOB     NULL,
			version INTEGER  NOT NULL DEFAULT 1
		)
	`, t.nodes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.nodes, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			from_id  INTEGER NOT NULL,
			type_id  INTEGER NOT NULL,
			to_id    INTEGER NOT NULL,
			sequence INTEGER NOT NULL,
			data     BLOB    NULL,

			PRIMARY KEY(from_id, type_id, to_id),
			FOREIGN KEY(from_id) REFERENCES %s(id) ON DELETE CASCADE
		)
	`, t.edges, t.nodes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_sequence ON %s (sequence)
	`, t.edges, t.edges)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s sequence index: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			from_id     INTEGER NOT NULL,
			type_id     INTEGER NOT NULL,
			count       INTEGER NOT NULL DEFAULT 0,
			version 	INTEGER NOT NULL NOT NULL DEFAULT 1,

			PRIMARY KEY(from_id, type_id),
			FOREIGN KEY(from_id) REFERENCES %s(id) ON DELETE CASCADE
		)
	`, t.edgeInfos, t.nodes)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s table: %w", t.edgeInfos, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS %s_insert AFTER INSERT ON %s 
		BEGIN
			INSERT INTO %s (from_id, type_id, count) VALUES (NEW.from_id, NEW.type_id, 1)
			ON CONFLICT(from_id, type_id) DO
				UPDATE SET count = count + 1, version = version + 1;
		END
	`, t.edges, t.edges, t.edgeInfos)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s insert trigger: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS %s_update AFTER UPDATE ON %s 
		BEGIN
			UPDATE %s SET version = version + 1
			WHERE from_id = NEW.from_id AND type_id = NEW.type_id;
		END
	`, t.edges, t.edges, t.edgeInfos)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s update trigger: %w", t.edges, err)
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS %s_insert AFTER DELETE ON %s 
		BEGIN
			UPDATE %s SET count = count - 1, version = version + 1
			WHERE from_id = OLD.from_id AND type_id = OLD.type_id;
		END
	`, t.edges, t.edges, t.edgeInfos)); err != nil {
		return tables{}, fmt.Errorf("failed to create %s delete trigger: %w", t.edges, err)
	}
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
