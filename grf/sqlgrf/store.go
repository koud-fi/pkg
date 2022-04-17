package sqlgrf

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/koud-fi/pkg/grf"
)

type store struct {
	db *sql.DB

	tablesMu  sync.Mutex
	tablesMap map[grf.NodeType]tables
}

func NewStore(db *sql.DB) grf.Store {
	return &store{db: db}
}

func (s *store) Node(nt grf.NodeType, id ...grf.LocalID) ([]grf.NodeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	if len(id) == 0 {
		return []grf.NodeData{}, nil
	}
	var ids strings.Builder
	for i, id := range id {
		if i > 0 {
			ids.WriteByte(',')
		}
		ids.WriteString(strconv.FormatInt(int64(id), 10))
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT id, data, ts FROM %s
		WHERE id IN (%s)
	`, t.nodes, idStr(id)))
	if err != nil {
		return nil, err
	}
	return scanNodes(rows, make([]grf.NodeData, 0, len(id)))
}

func (s *store) NodeRange(nt grf.NodeType, after grf.LocalID, limit int) ([]grf.NodeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT id, data, ts FROM %s
		WHERE id > %d
		ORDER BY id ASC
		LIMIT %d
	`, t.nodes, after, limit))
	if err != nil {
		return nil, err
	}
	return scanNodes(rows, make([]grf.NodeData, 0, limit))
}

func scanNodes(rows *sql.Rows, out []grf.NodeData) ([]grf.NodeData, error) {
	return scanRows(rows, out, func(rows *sql.Rows, nd *grf.NodeData) error {
		return rows.Scan(&nd.ID, &nd.Data, &nd.Timestamp)
	})
}

func (s *store) Edge(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeType, to ...grf.ID,
) ([]grf.EdgeData, error) {

	// ???

	panic("TODO")
}

func (s *store) EdgeCount(
	nt grf.NodeType, from grf.LocalID, et ...grf.EdgeType,
) (map[grf.EdgeType]int, error) {

	// ???

	panic("TODO")
}

func (s *store) EdgeRange(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeType, offset, limit int,
) ([]grf.EdgeData, error) {

	// ???

	panic("TODO")
}

func (s *store) AddNode(nt grf.NodeType, data []byte) (grf.LocalID, time.Time, error) {
	t, err := s.tables(nt)
	if err != nil {
		return 0, time.Time{}, err
	}
	ts := time.Now()
	res, err := s.db.Exec(fmt.Sprintf(`
		INSERT INTO %s (data, ts) VALUES (?, ?)
	`, t.nodes), data, ts)
	if err != nil {
		return 0, time.Time{}, err
	}
	id, err := res.LastInsertId()
	return grf.LocalID(id), ts, err
}

func (s *store) UpdateNode(nt grf.NodeType, id grf.LocalID, data []byte) error {
	t, err := s.tables(nt)
	if err != nil {
		return err
	}
	res, err := s.db.Exec(fmt.Sprintf(`
		UPDATE %s SET data = ?
	`, t.nodes), data)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return grf.ErrNotFound
	}
	return nil
}

func (s *store) DeleteNode(nt grf.NodeType, id ...grf.LocalID) error {
	t, err := s.tables(nt)
	if err != nil {
		return err
	}
	if len(id) == 0 {
		return nil
	}
	_, err = s.db.Exec(fmt.Sprintf(`
		DELETE FROM %s WHERE id IN (%s)
	`, t.nodes, idStr(id)))
	return err
}

func (s *store) SetEdge(nt grf.NodeType, e ...grf.EdgeData) error {

	// ???

	panic("TODO")
}

func (s *store) DeleteEdge(nt grf.NodeType, from grf.LocalID, et grf.EdgeType, to ...grf.ID) error {

	// ???

	panic("TODO")
}
