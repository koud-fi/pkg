package api

import "github.com/koud-fi/pkg/schema/jsonschema"

type (
	EndpointSchema struct {
		Name   string             `json:"name,omitempty" yaml:"name,omitempty"`
		Input  *jsonschema.Schema `json:"in,omitempty" yaml:"in,omitempty"`
		Output *jsonschema.Schema `json:"out,omitempty" yaml:"out,omitempty"`
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
		Input:  jsonschema.FromType(e.inType),
		Output: jsonschema.FromType(e.outType),
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
