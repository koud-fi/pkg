package api

import (
	"github.com/koud-fi/pkg/schema"
)

type (
	EndpointSchema struct {
		Input  schema.Type `json:"in,omitempty" yaml:"in,omitempty"`
		Output schema.Type `json:"out,omitempty" yaml:"out,omitempty"`
	}
	MuxSchema struct {
		Endpoints []MuxEndpointSchema `json:"endpoints" yaml:"endpoints"`
	}
	MuxEndpointSchema struct {
		Name string `json:"name" yaml:"name"`
		EndpointSchema
	}
)

func (e *Endpoint) Schema() EndpointSchema {

	// TODO: this can't really handle all possible output types

	return EndpointSchema{
		Input:  schema.ResolveTypeOf(e.inType),
		Output: schema.ResolveTypeOf(e.outType),
	}
}

func (m *Mux) Schema() MuxSchema {
	schema := MuxSchema{
		Endpoints: make([]MuxEndpointSchema, 0, len(m.endpoints)),
	}
	for name, endpoint := range m.endpoints {
		schema.Endpoints = append(schema.Endpoints, MuxEndpointSchema{
			Name:           name,
			EndpointSchema: endpoint.Schema(),
		})
	}
	return schema
}
