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
	Count   int   `json:"count"`
	Version int64 `json:"version"`
}

func (e Edge[T]) String() string {
	return fmt.Sprintf("%d>%s>%d(%d) %v", e.From, e.Type, e.To, e.Sequence, e.Data)
}

func LookupEdge[T any](g *Graph, from ID, et EdgeType, to ID) (*Edge[T], error) {

	// ???

	panic("TODO")
}

func LookupEdgeInfo(g *Graph, from ID, et ...EdgeType) (map[EdgeType]EdgeInfo, error) {
	/*
		ti, s, err := g.parseID(from)
		if err != nil {
			return nil, err
		}
		etIDs := make([]EdgeTypeID, len(et))
		for i := range et {
			var ok bool
			if etIDs[i], ok = ti.edgeTypeMap[et[i]]; !ok {
				return nil, ErrInvalidEdgeType
			}
		}
		m, err := s.EdgeInfo(ti.Type, from.localID(), etIDs...)
		if err != nil {
			return nil, err
		}
	*/

	// ???

	panic("TODO")
}

func EdgeRange[T any](g *Graph, from ID, et EdgeType, offset, limit int) ([]Edge[T], error) {

	// ???

	panic("TODO")
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
	ti, etID, s, err := g.parseEdgeParams(from, et)
	if err != nil {
		return err
	}
	return s.DeleteEdge(ti.Type, from.localID(), etID, to...)
}
