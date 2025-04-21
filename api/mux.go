package api

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/assign"
)

type Mux struct {
	endpoints map[string]*Endpoint
	converter *assign.Converter
}

func NewMux() *Mux {
	return &Mux{
		endpoints: make(map[string]*Endpoint),
		converter: assign.NewDefaultConverter(),
	}
}

func (m *Mux) Register(name string, fn any) {
	e, err := NewEndpoint(m.converter, fn)
	if err != nil {
		panic(err)
	}
	m.endpoints[name] = e
}

func (m *Mux) Call(ctx context.Context, name string, args Arguments) (any, error) {
	e, ok := m.endpoints[name]
	if !ok {
		return nil, fmt.Errorf("endpoint %q not found", name)
	}
	return e.Call(ctx, args)
}
