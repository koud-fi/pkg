package grf

import (
	"encoding/json"
	"fmt"
	"time"
)

type Node struct {
	id ID
	t  NodeType
	d  NodeData
}

type NodeType string

func (n Node) ID() ID         { return n.id }
func (n Node) Type() NodeType { return n.t }

func (n Node) Unmarshal(v any) error {
	return json.Unmarshal(n.d.Data, v)
}

func (n Node) Timestamp() time.Time { return n.d.Timestamp }

func (n Node) String() string {
	return fmt.Sprintf("%d(%s)(%v) %s",
		n.id, n.t, n.d.Timestamp.UTC().Format(time.RFC3339Nano), string(n.d.Data))
}
