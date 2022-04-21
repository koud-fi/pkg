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

func (e Edge[T]) String() string {
	return fmt.Sprintf("%d>%s>%d(%d) %v", e.From, e.Type, e.To, e.Sequence, e.Data)
}

// TODO: edge lookup / listing

func (g *Graph) SetEdge(e ...Edge[any]) error {
	ed := make([]struct {
		EdgeData
		ti TypeInfo
		s  Store
	}, len(e))
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
			ed[i].Sequence = time.Now().UnixNano()
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

// TODO: edge deletion
