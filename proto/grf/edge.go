package grf

import (
	"fmt"
	"time"
)

type Edge[T any] struct {
	From     ID
	Type     EdgeType
	To       ID
	Sequence int64
	Data     T
}

type EdgeType string

type EdgeInfo struct {
	Count   int
	Version int64
}

func (e Edge[T]) String() string {
	return fmt.Sprintf("%d>%s>%d(%d) %v", e.From, e.Type, e.To, e.Sequence, e.Data)
}

func LookupEdge[T any](g *Graph, from ID, et EdgeType, to ID) (*Edge[T], error) {
	ti, etID, s, err := g.parseEdgeArgs(from, et)
	if err != nil {
		return nil, err
	}
	eds, err := s.Edge(ti.Type, from.localID(), etID, to)
	if err != nil {
		return nil, err
	}
	if len(eds) == 0 {
		return nil, ErrNotFound
	}
	es, err := convertEdges[T](from, et, eds)
	if err != nil {
		return nil, err
	}
	return &es[0], nil
}

func LookupEdgeInfo(g *Graph, from ID, et ...EdgeType) (map[EdgeType]EdgeInfo, error) {
	ti, s, err := g.parseID(from)
	if err != nil {
		return nil, err
	}
	var (
		etIDs = make([]EdgeTypeID, len(et))
		etMap = make(map[EdgeTypeID]EdgeType, len(et))
	)
	for i := range et {
		var ok bool
		if etIDs[i], ok = ti.edgeTypeMap[et[i]]; !ok {
			return nil, ErrInvalidEdgeType
		}
		etMap[etIDs[i]] = et[i]
	}
	eids, err := s.EdgeInfo(ti.Type, from.localID(), etIDs...)
	if err != nil {
		return nil, err
	}
	m := make(map[EdgeType]EdgeInfo, len(eids))
	for _, eid := range eids {
		m[etMap[eid.TypeID]] = eid.EdgeInfo
	}
	return m, nil
}

func EdgeRange[T any](g *Graph, from ID, et EdgeType, offset, limit int) ([]Edge[T], error) {
	ti, etID, s, err := g.parseEdgeArgs(from, et)
	if err != nil {
		return nil, err
	}
	eds, err := s.EdgeRange(ti.Type, from.localID(), etID, offset, limit)
	if err != nil {
		return nil, err
	}
	return convertEdges[T](from, et, eds)
}

func convertEdges[T any](from ID, et EdgeType, eds []EdgeData) ([]Edge[T], error) {
	es := make([]Edge[T], 0, len(eds))
	for _, ed := range eds {
		v, err := unmarshal[T](nil, ed.Data)
		if err != nil {
			return nil, fmt.Errorf("invalid edge data: %w", err)
		}
		es = append(es, Edge[T]{
			From:     from,
			Type:     et,
			To:       ed.To,
			Sequence: ed.Sequence,
			Data:     v,
		})
	}
	return es, nil
}

func SetEdge(g *Graph, e ...Edge[any]) error {
	var (
		ed = make([]struct {
			EdgeData
			ti TypeInfo
			s  Store
		}, len(e))
		seq = time.Now().UnixNano()
	)
	for i := range e {
		var err error
		if ed[i].ti, ed[i].s, err = g.parseID(e[i].From); err != nil {
			return err
		}
		ed[i].From = e[i].From.localID()
		ed[i].To = e[i].To
		ed[i].Sequence = e[i].Sequence
		ed[i].Data = marshal(e[i].Data)

		if etID, ok := ed[i].ti.edgeTypeMap[e[i].Type]; ok {
			ed[i].TypeID = etID
		} else {
			return ErrInvalidEdgeType
		}
		if ed[i].Sequence == 0 {
			ed[i].Sequence = seq
			seq++
		}
	}
	for _, ed := range ed {

		// TODO: batch by store/type

		if err := ed.s.SetEdge(ed.ti.Type, ed.EdgeData); err != nil {
			return err
		}
	}
	return nil
}

func DeleteEdge(g *Graph, from ID, et EdgeType, to ...ID) error {
	ti, etID, s, err := g.parseEdgeArgs(from, et)
	if err != nil {
		return err
	}
	return s.DeleteEdge(ti.Type, from.localID(), etID, to...)
}
