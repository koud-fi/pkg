package schema

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

const (
	String  TypeName = "string"
	Number  TypeName = "number"
	Integer TypeName = "integer"
	Object  TypeName = "object"
	Array   TypeName = "array"
	Boolean TypeName = "boolean"
	Null    TypeName = "null"

	// TODO: Date     Format = "date"
	DateTime Format = "date-time"
	// TODO: Password Format = "password"
	// TODO: Byte     Format = "byte" // base64-encoded characters
	// TODO: Binary   Format = "binary"
)

type TypeName string
type Format string

type Type struct {
	Type       TypeName   `json:"type,omitempty"`
	Format     Format     `json:"format,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Items      *Type      `json:"items,omitempty"`

	Tags map[string]string `json:"tags,omitempty"`
}

// ResolveType infers a JSON-Schema Type for a generic Go type T.
func ResolveType[T any](opts ...Option) Type {
	var v T
	return ResolveTypeOf(v, opts...)
}

// ResolveTypeOf infers a JSON-Schema Type for the value or type v.
func ResolveTypeOf(v any, opts ...Option) Type {
	c := config{
		customTypes: map[typeKey]func(reflect.Type) Type{
			{"time", "Time"}: func(_ reflect.Type) Type {
				return Type{Type: String, Format: DateTime}
			},
		},
		tags:       []string{},
		inProgress: make(map[reflect.Type]*Type),
	}
	for _, opt := range opts {
		opt(&c)
	}
	return resolveType(c, v)
}

// JSONName returns the JSON field name based on the struct tag, if present.
func (t Type) JSONName(key string) string {
	if jsonTag, ok := t.Tags["json"]; ok {
		return strings.Split(jsonTag, ",")[0]
	}
	return key
}

// ExampleValue produces a native Go value matching the schema.
func (t Type) ExampleValue() any {
	switch t.Type {
	case String:
		if t.Format == DateTime {
			return time.RFC3339
		}
		return ""
	case Number, Integer:
		return 0
	case Object:
		o := make(map[string]any, len(t.Properties))
		for name, typ := range t.Properties {
			o[name] = typ.ExampleValue()
		}
		return o
	case Array:
		if t.Items == nil {
			return []any{nil}
		}
		return []any{t.Items.ExampleValue()}
	case Boolean:
		return false
	default:
		return nil
	}
}

// ExampleJSON returns a pretty-printed JSON example for the schema.
func (t Type) ExampleJSON() string {
	data, _ := json.MarshalIndent(t.ExampleValue(), "", "\t")
	return string(data)
}

// resolveType is the core recursive resolver with cycle detection.
func resolveType(c config, v any) Type {
	// Handle map[string]any as an object literal
	switch vv := v.(type) {
	case map[string]any:
		full := Type{Type: Object}
		full.allocProps(len(vv))
		full.Properties.fromMap(c, vv)
		return full
	case []any:
		full := Type{Type: Array}
		for _, elem := range vv {
			it := resolveType(c, elem)
			full.Items = &it
		}
		return full
	}

	// Fallback to reflect
	var rt reflect.Type
	switch vv := v.(type) {
	case reflect.Type:
		rt = vv
	default:
		rt = reflect.TypeOf(v)
	}
	if rt == nil {
		return Type{}
	}

	// Cycle detection: if we already started resolving this struct type, emit a $ref.
	var placeholder *Type
	if rt.Kind() == reflect.Struct {
		if _, seen := c.inProgress[rt]; seen {
			// Emit a $ref to break the cycle
			return Type{Tags: map[string]string{"$ref": "#/definitions/" + rt.Name()}}
		}
		placeholder = &Type{}
		c.inProgress[rt] = placeholder
	}

	// Custom type overrides (e.g. time.Time)
	if ct, ok := c.customTypes[typeKey{rt.PkgPath(), strings.TrimLeft(rt.Name(), "*")}]; ok {
		full := ct(rt)
		*placeholder = full
		return full
	}

	// Build the full schema for rt
	var full Type
	switch rt.Kind() {
	case reflect.Ptr:
		full = resolveType(c, rt.Elem())

	case reflect.String:
		full.Type = String

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		full.Type = Integer

	case reflect.Float32, reflect.Float64:
		full.Type = Number

	case reflect.Bool:
		full.Type = Boolean

	case reflect.Struct:
		full.Type = Object
		full.allocProps(rt.NumField())
		full.Properties.fromStructFields(c, rt)

	case reflect.Map:
		full.Type = Object // free-form map

	case reflect.Slice:
		item := resolveType(c, rt.Elem())
		full.Type = Array
		full.Items = &item

	case reflect.Interface:
		// no constraints for interface

	default:
		panic("cannot resolve schema for type: " + rt.Kind().String())
	}

	// If we registered a placeholder, fill it now and return
	if placeholder != nil {
		*placeholder = full
	}
	return full
}

// allocProps initializes the Properties map if nil.
func (t *Type) allocProps(lenHint int) {
	if t.Properties == nil {
		t.Properties = make(Properties, lenHint)
	}
}
