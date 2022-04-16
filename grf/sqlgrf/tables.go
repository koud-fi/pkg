package sqlgrf

import (
	"fmt"

	"github.com/koud-fi/pkg/grf"
)

const (
	nodeTableSuffix     = "_node"
	edgeTypeTableSuffix = "_edgetype"
	edgeTableSuffix     = "_edge"
	keymapTableSuffix   = "_keymap"
)

type tables struct {
	nodes     string
	edgeTypes string
	edges     string
}

func (s *store) tables(nt grf.NodeType) (tables, error) {
	s.tablesMu.Lock()
	defer s.tablesMu.Unlock()

	if t, ok := s.tablesMap[nt]; ok {
		return t, nil
	}
	t := tables{
		nodes:     string(nt) + nodeTableSuffix,
		edgeTypes: string(nt) + edgeTypeTableSuffix,
		edges:     string(nt) + edgeTableSuffix,
	}
	if _, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB,
			ts   DATETIME NOT NULL
		)
	`, t.nodes)); err != nil {
		return tables{}, err
	}

	// TODO: create edgetype table
	// TODO: create edge table

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
