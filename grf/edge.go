package grf

import (
	"encoding/json"
	"fmt"
)

type Edge struct {
	from ID
	d    EdgeData
}

type EdgeType string

func NewEdge(from ID, et EdgeType, to ID, seq int64, v any) Edge {
	return Edge{
		from: from,
		d:    EdgeData{Type: et, To: to, Sequence: seq, Data: marshal(v)},
	}
}

func (e Edge) From() ID        { return e.from }
func (e Edge) Type() EdgeType  { return e.d.Type }
func (e Edge) To() ID          { return e.d.To }
func (e Edge) Sequence() int64 { return e.d.Sequence }

func (e Edge) Unmarshal(v any) error {
	return json.Unmarshal(e.d.Data, v)
}

func (e Edge) String() string {
	return fmt.Sprintf("%d>%s>%d(%d) %s",
		e.from, e.d.Type, e.d.To, e.d.Sequence, string(e.d.Data))
}
