package grf

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Node struct {
	id  ID
	d   NodeData
	ti  TypeInfo
	err error
}

type NodeType string

func (n Node) ID() ID         { return n.id }
func (n Node) Type() NodeType { return n.ti.Type }

func (n Node) Data() (any, error) {
	if n.err != nil {
		return nil, n.err
	}
	var v any
	if n.ti.dataType != nil {
		v = reflect.New(n.ti.dataType).Elem().Interface()
	}
	return v, json.Unmarshal(n.d.Data, &v)
}

func (n Node) Timestamp() time.Time { return n.d.Timestamp }

func (n Node) String() string {
	ts := n.d.Timestamp.UTC().Format(time.RFC3339Nano)
	return fmt.Sprintf("%d(%s)(%v) %s", n.id, n.ti.Type, ts, string(n.d.Data))
}
