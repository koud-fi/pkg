package sqlgrf

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

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
		return rows.Scan(&nd.ID, &nd.Data, &nd.Timestamp)
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
		SELECT from_id, type_id, to_id, seq, data FROM %s
		WHERE from_id = ? AND type_id = ? AND to_id IN (%s)
	`, t.edges, idStr(to)), from, et)
	if err != nil {
		return nil, err
	}
	return scanEdges(rows, make([]grf.EdgeData, 0, len(to)))
}

func (s *store) EdgeInfo(
	nt grf.NodeType, from grf.LocalID, et ...grf.EdgeTypeID,
) (map[grf.EdgeTypeID]grf.EdgeInfo, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	if len(et) == 0 {
		return make(map[grf.EdgeTypeID]grf.EdgeInfo), nil
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT type_id, count, last_update
		FROM %s WHERE from_id = ? AND type_id IN (%s)
	`, t.edgeInfos, idStr(et)), from)
	if err != nil {
		return nil, err
	}
	out := make(map[grf.EdgeTypeID]grf.EdgeInfo, len(et))
	for rows.Next() {
		var (
			et   grf.EdgeTypeID
			info grf.EdgeInfo
		)
		if err := rows.Scan(&et, &info.Count, &info.LastUpdate); err != nil {
			return nil, err
		}
		out[et] = info
	}
	return out, nil
}

func (s *store) EdgeRange(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, offset, limit int,
) ([]grf.EdgeData, error) {
	t, err := s.tables(nt)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(fmt.Sprintf(`
		SELECT from_id, type_id, to_id, seq, data FROM %s
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
		UPDATE %s SET data = ? WHERE id = ?
	`, t.nodes), data, id)
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
