package sqlgrf

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/koud-fi/pkg/grf"
)

type store struct {
	db *sql.DB

	tablesMu    sync.Mutex
	tablesMap   map[grf.NodeType]tables
	edgeTypeMu  sync.Mutex
	edgeTypeMap map[[2]string]int64
}

func NewStore(db *sql.DB) grf.Store {
	return &store{db: db, edgeTypeMap: make(map[[2]string]int64)}
}

func (s *store) Node(nt grf.NodeType, id ...grf.LocalID) ([]grf.NodeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	if len(id) == 0 {
		return []grf.NodeData{}, nil
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT id, data, version FROM %s
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
		SELECT id, data, version FROM %s
		WHERE id > ?
		ORDER BY id ASC
		LIMIT ?
	`, t.nodes), after, limit)
	if err != nil {
		return nil, err
	}
	return scanNodes(rows, make([]grf.NodeData, 0, min(limit, 50)))
}

func scanNodes(rows *sql.Rows, out []grf.NodeData) ([]grf.NodeData, error) {
	return scanRows(rows, out, func(rows *sql.Rows, nd *grf.NodeData) error {
		return rows.Scan(&nd.ID, &nd.Data, &nd.Version)
	})
}

func (s *store) Edge(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, to ...grf.ID,
) ([]grf.EdgeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	if len(to) == 0 {
		return []grf.EdgeData{}, nil
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT from_id, type_id, to_id, sequence, data FROM %s
		WHERE from_id = ? AND type_id = ? AND to_id IN (%s)
	`, t.edges, idStr(to)), from, et)
	if err != nil {
		return nil, err
	}
	return scanEdges(rows, make([]grf.EdgeData, 0, len(to)))
}

func (s *store) EdgeInfo(
	nt grf.NodeType, from grf.LocalID, et ...grf.EdgeTypeID,
) ([]grf.EdgeInfoData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	if len(et) == 0 {
		return []grf.EdgeInfoData{}, nil
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT type_id, count, version
		FROM %s WHERE from_id = ? AND type_id IN (%s)
	`, t.edgeInfos, idStr(et)), from)
	if err != nil {
		return nil, err
	}
	out := make([]grf.EdgeInfoData, 0, len(et))
	return scanRows(rows, out, func(rows *sql.Rows, eid *grf.EdgeInfoData) error {
		return rows.Scan(&eid.TypeID, &eid.Count, &eid.Version)
	})
}

func (s *store) EdgeRange(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, offset, limit int,
) ([]grf.EdgeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT from_id, type_id, to_id, sequence, data FROM %s
		WHERE from_id = ? AND type_id = ?
		ORDER BY seq DESC
		LIMIT ? OFFSET ?
	`, t.edges), from, et, limit, offset)
	if err != nil {
		return nil, err
	}
	return scanEdges(rows, make([]grf.EdgeData, 0, min(limit, 50)))
}

func scanEdges(rows *sql.Rows, out []grf.EdgeData) ([]grf.EdgeData, error) {
	return scanRows(rows, out, func(rows *sql.Rows, ed *grf.EdgeData) error {
		return rows.Scan(&ed.From, &ed.TypeID, &ed.To, &ed.Sequence, &ed.Data)
	})
}

func (s *store) AddNode(nt grf.NodeType, data []byte) (grf.LocalID, int64, error) {
	t, err := s.tables(nt)
	if err != nil {
		return 0, 0, err
	}
	ver := int64(1)
	res, err := s.db.Exec(fmt.Sprintf(`
		INSERT INTO %s (data, version) VALUES (?, ?)
	`, t.nodes), data, ver)
	if err != nil {
		return 0, 0, err
	}
	id, err := res.LastInsertId()
	return grf.LocalID(id), ver, err
}

func (s *store) UpdateNode(
	nt grf.NodeType, id grf.LocalID, data []byte, currentVersion int64,
) error {
	t, err := s.tables(nt)
	if err != nil {
		return err
	}
	res, err := s.db.Exec(fmt.Sprintf(`
		UPDATE %s SET data = ?, version = version + 1
		WHERE id = ? AND version = ?
	`, t.nodes), data, id, currentVersion)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		var exists bool
		if err := s.db.QueryRow(fmt.Sprintf(`
			SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)
		`, t.nodes), id).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return grf.ErrVersionMismatch
		}
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
	t, err := s.tables(nt)
	if err != nil {
		return err
	}
	for _, e := range e {
		if _, err := s.db.Exec(fmt.Sprintf(`
			INSERT INTO %s (from_id, type_id, to_id, sequence, data) VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(from_id, type_id, to_id) DO
				UPDATE SET sequence = excluded.sequence, data = excluded.data
		`, t.edges), e.From, e.TypeID, e.To, e.Sequence, e.Data); err != nil {
			return err
		}
	}
	return nil
}

func (s *store) DeleteEdge(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, to ...grf.ID,
) error {
	t, err := s.tables(nt)
	if err != nil {
		return err
	}
	if len(to) == 0 {
		return nil
	}
	_, err = s.db.Exec(fmt.Sprintf(`
		DELETE FROM %s
		WHERE from_id = ? AND type_id = ? AND to_id IN (%s)
	`, t.edges, idStr(to)), from, et)
	return err
}
