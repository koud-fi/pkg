package memgrf

import (
	"sort"

	"github.com/koud-fi/pkg/grf"
)

type store struct {
	locker
	Data map[grf.NodeType]nodeList `json:"data"`
}

type nodeList struct {
	LastInsertID grf.LocalID `json:"lastInsertId,omitempty"`
	Nodes        []node      `json:"nodes"`
}

type node struct {
	grf.NodeData
	Edges     map[grf.EdgeType][]grf.EdgeData `json:"edges"`
	IsDeleted bool                            `json:"isDeleted,omitempty"`
}

func NewStore() grf.Store {
	return &store{Data: make(map[grf.NodeType]nodeList)}
}

func (s *store) Node(nt grf.NodeType, id ...grf.LocalID) ([]grf.NodeData, error) {
	defer s.rlock()()
	var (
		nl  = s.Data[nt]
		out = make([]grf.NodeData, 0, len(id))
	)
	for _, id := range id {
		if i, ok := searchNode(nl.Nodes, id); ok {
			if !nl.Nodes[i].IsDeleted {
				out = append(out, nl.Nodes[i].NodeData)
			}
		}
	}
	return out, nil
}

func (s *store) NodeRange(
	nt grf.NodeType, after grf.LocalID, limit int,
) ([]grf.NodeData, error) {
	defer s.rlock()()
	var (
		nl       = s.Data[nt]
		start, _ = searchNode(nl.Nodes, after)
		ns       = nl.Nodes[start:]
	)
	if len(ns) < limit {
		limit = len(ns)
	}
	out := make([]grf.NodeData, 0, limit)
	for _, n := range ns {
		if n.IsDeleted {
			continue
		}
		out = append(out, n.NodeData)
		if len(out) == limit {
			break
		}
	}
	return out, nil
}

func (s *store) Edge(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, to ...grf.ID,
) ([]grf.EdgeData, error) {
	defer s.rlock()()

	// ???

	panic("TODO")
}

func (s *store) EdgeInfo(
	nt grf.NodeType, from grf.LocalID, et ...grf.EdgeTypeID,
) (map[grf.EdgeTypeID]grf.EdgeInfo, error) {
	defer s.rlock()()

	// ???

	panic("TODO")
}

func (s *store) EdgeRange(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, offset, limit int,
) ([]grf.EdgeData, error) {
	defer s.rlock()()

	// ???

	panic("TODO")
}

func (s *store) AddNode(nt grf.NodeType, data []byte) (grf.LocalID, int64, error) {
	defer s.lock()()
	nl := s.Data[nt]
	nl.LastInsertID++
	ver := int64(1)
	nl.Nodes = append(nl.Nodes, node{
		NodeData: grf.NodeData{
			ID:      nl.LastInsertID,
			Data:    data,
			Version: ver,
		},
		Edges: make(map[grf.EdgeType][]grf.EdgeData),
	})
	s.Data[nt] = nl
	return nl.LastInsertID, ver, nil
}

func (s *store) UpdateNode(
	nt grf.NodeType, id grf.LocalID, data []byte, currentVersion int64,
) error {
	defer s.lock()()
	var (
		nl    = s.Data[nt]
		i, ok = searchNode(nl.Nodes, id)
	)
	if !ok || nl.Nodes[i].IsDeleted {
		return grf.ErrNotFound
	}
	if nl.Nodes[i].Version != currentVersion {
		return grf.ErrVersionMismatch
	}
	nl.Nodes[i].Data = data
	nl.Nodes[i].Version++
	return nil
}

func (s *store) DeleteNode(nt grf.NodeType, id ...grf.LocalID) error {
	defer s.lock()()
	nl := s.Data[nt]
	for _, id := range id {
		if i, ok := searchNode(nl.Nodes, id); ok {
			n := &nl.Nodes[i]
			n.Data = nil
			n.Edges = nil
			n.IsDeleted = true
		}
	}
	return nil
}

func (s *store) SetEdge(nt grf.NodeType, e ...grf.EdgeData) error {
	defer s.lock()()

	// ???

	panic("TODO")
}

func (s *store) DeleteEdge(
	nt grf.NodeType, from grf.LocalID, et grf.EdgeTypeID, to ...grf.ID,
) error {
	defer s.lock()()

	// ???

	panic("TODO")
}

func searchNode(ns []node, id grf.LocalID) (int, bool) {
	i := sort.Search(len(ns), func(i int) bool {
		return ns[i].ID >= id
	})
	return i, i < len(ns) && ns[i].ID == id
}
