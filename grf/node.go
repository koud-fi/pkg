package grf

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Node struct {
	s   Store
	id  ID
	d   NodeData
	ti  TypeInfo
	err error
}

type NodeType string

func (n Node) ID() ID         { return n.id }
func (n Node) Type() NodeType { return n.ti.Type }

func (n Node) Data() (v any, err error) {
	if n.err != nil {
		return nil, n.err
	}
	if n.ti.dataType != nil {
		p := reflect.New(n.ti.dataType)
		if n.d.Data != nil {
			err = json.Unmarshal(n.d.Data, p.Interface())
		}
		v = p.Elem().Interface()
	} else if n.d.Data != nil {
		err = json.Unmarshal(n.d.Data, &v)
	}
	return
}

func (n Node) Timestamp() time.Time { return n.d.Timestamp }

func (n Node) String() string {
	ts := n.d.Timestamp.UTC().Format(time.RFC3339Nano)
	return fmt.Sprintf("%d(%s)(%v) %s", n.id, n.ti.Type, ts, string(n.d.Data))
}

func (n *Node) Update(fn func(v any) (any, error)) *Node {
	if n.err != nil {
		return n
	}
	var v any
	if v, n.err = n.Data(); n.err != nil {
		return n
	}
	if v, n.err = fn(v); n.err != nil {
		return n
	}
	n.d.Data = marshal(v)
	n.err = n.s.UpdateNode(n.ti.Type, n.d.ID, n.d.Data)
	return n
}

func (n Node) Err() error { return n.err }
