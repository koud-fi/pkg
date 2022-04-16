package grf

import "encoding/json"

type Edge struct {
	d EdgeData
}

type EdgeType string

func (e Edge) Data(v any) error {
	return json.Unmarshal(e.d.Data, v)
}
