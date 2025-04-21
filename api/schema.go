package api

import "github.com/koud-fi/pkg/schema/jsonschema"

type EndpointSchema struct {
	Input  *jsonschema.Schema `json:"input,omitempty" yaml:"input,omitempty"`
	Output *jsonschema.Schema `json:"output,omitempty" yaml:"output,omitempty"`
}

func (e Endpoint) Schema() *EndpointSchema {

	// TODO: this can't really handle all possible output types

	return &EndpointSchema{
		Input:  jsonschema.FromType(e.inType),
		Output: jsonschema.FromType(e.outType),
	}
}
