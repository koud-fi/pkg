package jsonschema

import (
	"reflect"
	"strings"

	"github.com/koud-fi/pkg/schema/flex"
)

const (
	Boolean Type = "boolean"
	Integer Type = "integer"
	Number  Type = "number"
	String  Type = "string"
	Array   Type = "array"
	Object  Type = "object"
	Null    Type = "null"
)

// Schema models a JSON Schema object and can be encoded/decoded as JSON or YAML.
type (
	Schema struct {
		Schema               string               `json:"$schema,omitempty" yaml:"$schema,omitempty"`
		Title                string               `json:"title,omitempty" yaml:"title,omitempty"`
		Description          string               `json:"description,omitempty" yaml:"description,omitempty"`
		Type                 flex.OneOrMany[Type] `json:"type,omitempty" yaml:"type,omitempty"`
		Properties           map[string]*Schema   `json:"properties,omitempty" yaml:"properties,omitempty"`
		Items                *Schema              `json:"items,omitempty" yaml:"items,omitempty"`
		Required             []string             `json:"required,omitempty" yaml:"required,omitempty"`
		AdditionalProperties interface{}          `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
		// Extend with fields like Enum, Format, etc.
	}
	Type string
)

// FromType inspects a Go type and returns the corresponding JSON Schema.
func FromType(t reflect.Type) *Schema {
	// Dereference pointers
	if t.Kind() == reflect.Ptr {
		return FromType(t.Elem())
	}
	var s Schema
	switch t.Kind() {
	case reflect.Bool:
		s.Type = flex.One(Boolean)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = flex.One(Integer)

	case reflect.Float32, reflect.Float64:
		s.Type = flex.One(Number)

	case reflect.String:
		s.Type = flex.One(String)

	case reflect.Slice, reflect.Array:
		s.Type = flex.One(Array)
		s.Items = FromType(t.Elem())

	case reflect.Map:
		s.Type = flex.One(Object)
		s.AdditionalProperties = FromType(t.Elem())

	case reflect.Struct:
		s.Type = flex.One(Object)
		s.Properties = make(map[string]*Schema)

		var required []string
		for i := range t.NumField() {
			f := t.Field(i)

			// skip unexported
			if f.PkgPath != "" {
				continue
			}

			// JSON field name
			name := f.Tag.Get("json")
			if name == "" || name == "-" {
				name = f.Name
			} else if idx := strings.Index(name, ","); idx != -1 {
				name = name[:idx]
			}
			child := FromType(f.Type)

			// schema tag for description
			desc := f.Tag.Get("schema")
			if desc != "" {
				child.Description = desc
			}
			s.Properties[name] = child

			// mark required if "omitempty" is not present
			if !strings.Contains(f.Tag.Get("json"), "omitempty") {
				required = append(required, name)
			}
		}
		if len(required) > 0 {
			s.Required = required
		}
	default:
		s.Type = flex.One(String)
	}
	return &s
}
