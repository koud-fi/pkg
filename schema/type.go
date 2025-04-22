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

type Type struct {
	Type       TypeName   `json:"type,omitempty"`
	Format     Format     `json:"format,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Items      *Type      `json:"items,omitempty"`

	Tags map[string]string `json:"tags,omitempty"`
}

type (
	TypeName string
	Format   string
)

func ResolveType[T any](opts ...Option) Type {
	var v T
	return ResolveTypeOf(v, opts...)
}

func ResolveTypeOf(v any, opts ...Option) Type {
	c := config{
		customTypes: map[typeKey]func(reflect.Type) Type{
			{"time", "Time"}: func(_ reflect.Type) Type {
				return Type{Type: String, Format: DateTime}
			},
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return Type{}.resolve(c, v)
}

func (t Type) JSONName(key string) string {
	if jsonTag, ok := t.Tags["json"]; ok {
		return strings.Split(jsonTag, ",")[0]
	}
	return key // TODO: proper JSON formatting
}

func (t Type) ExampleValue() any {
	switch t.Type {
	case String:
		switch t.Format {
		case DateTime:
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

func (t Type) ExampleJSON() string {
	data, _ := json.MarshalIndent(t.ExampleValue(), "", "\t")
	return string(data)
}

func (t Type) resolve(c config, v any) Type {
	var rt reflect.Type
	switch v := v.(type) {
	case reflect.Type:
		rt = v
	case map[string]any:
		t.Type = Object
		t.allocProps(len(v))
		t.Properties.fromMap(c, v)

	case []any:
		t.Type = Array
		for i := range v {
			if t.Items == nil {
				t.Items = new(Type)
			}

			// TODO: handle possible mixed types

			it := t.Items.resolve(c, v[i])
			t.Items = &it
		}
	default:
		rt = reflect.TypeOf(v)
	}
	if rt == nil {
		return t
	}
	if ct, ok := c.customTypes[typeKey{rt.PkgPath(), strings.TrimLeft(rt.Name(), "*")}]; ok {
		return ct(rt)
	}
	switch rt.Kind() {
	case reflect.Ptr:
		t = Type{}.resolve(c, rt.Elem())

	case reflect.String:
		t.Type = String

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
		reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:

		t.Type = Integer

	case reflect.Float64, reflect.Float32:
		t.Type = Number

	case reflect.Bool:
		t.Type = Boolean

	case reflect.Struct:
		t.Type = Object
		t.allocProps(rt.NumField())
		t.Properties.fromStructFields(c, rt)

	case reflect.Map:
		t.Type = Object

		// TODO: missing fields

	case reflect.Slice:
		t.Type = Array
		it := Type{}.resolve(c, rt.Elem())
		t.Items = &it

	case reflect.Interface:
		// TODO: ???

	default:
		panic("cannot resolve schema for type: " + rt.Kind().String())
	}
	return t
}

func (t *Type) allocProps(lenHint int) {
	if t.Properties == nil {
		t.Properties = make(map[string]Type, lenHint)
	}
}
