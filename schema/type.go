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
	Object  TypeName = "object"
	Array   TypeName = "array"
	Boolean TypeName = "boolean"
	Null    TypeName = "null"

	DateTime Format = "date-time"
)

type Type struct {
	Type       TypeName   `json:"type,omitempty"`
	Format     Format     `json:"format,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Items      *Type      `json:"items,omitempty"`

	Tags map[string]string `json:"tags,omitempty"` // TODO
}

type (
	TypeName string
	Format   string
)

func (t Type) ExampleValue() any {
	switch t.Type {
	case String:
		switch t.Format {
		case DateTime:
			return time.RFC3339
		}
		return ""

	case Number:
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

type Properties map[string]Type

func (p Properties) fromStructFields(c config, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.Anonymous {
			p.fromStructFields(c, sf.Type)
			continue
		}
		name := sf.Name
		if jsonTag, ok := sf.Tag.Lookup("json"); ok {
			if name = strings.Split(jsonTag, ",")[0]; name == "-" {
				continue
			}
		}
		p[name] = resolveType(c, sf.Type)
	}
}

func resolveType(c config, rt reflect.Type) (t Type) {
	if ct, ok := c.customTypes[typeKey{rt.PkgPath(), strings.TrimLeft(rt.Name(), "*")}]; ok {
		return ct(rt)
	}
	switch rt.Kind() {
	case reflect.Ptr:
		t = resolveType(c, rt.Elem())

	case reflect.String:
		t.Type = String

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
		reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8,
		reflect.Float64, reflect.Float32:

		t.Type = Number

	case reflect.Struct:
		t.Type = Object
		t.Properties = make(map[string]Type, rt.NumField())
		t.Properties.fromStructFields(c, rt)

	case reflect.Slice:
		t.Type = Array
		it := resolveType(c, rt.Elem())
		t.Items = &it

	case reflect.Bool:
		t.Type = Boolean

	default:
		panic("cannot resolve schema for type: " + rt.Kind().String())
	}
	return
}
