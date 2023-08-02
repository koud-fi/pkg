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

	Tags map[string]string `json:"tags,omitempty"`
}

type (
	TypeName string
	Format   string
)

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

func (t *Type) allocProps(lenHint int) {
	if t.Properties == nil {
		t.Properties = make(map[string]Type, lenHint)
	}
}

type Properties map[string]Type

func (p Properties) fromStructFields(c config, rt reflect.Type) {
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if sf.Anonymous {
			p.fromStructFields(c, sf.Type)
			continue
		}
		t := Type{}.resolve(c, sf.Type)
		for _, tag := range c.tags {
			if v, ok := sf.Tag.Lookup(tag); ok {
				if t.Tags == nil {
					t.Tags = map[string]string{tag: v}
				} else {
					t.Tags[tag] = v
				}
			}
		}
		p[sf.Name] = t
	}
}

func (p Properties) fromMap(c config, m map[string]any) {
	for k, v := range m {

		// TODO: handle possible mixed types for a same field

		p[k] = Type{}.resolve(c, v)
	}
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
		reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8,
		reflect.Float64, reflect.Float32:

		t.Type = Number

	case reflect.Bool:
		t.Type = Boolean

	case reflect.Struct:
		t.Type = Object
		t.allocProps(rt.NumField())
		t.Properties.fromStructFields(c, rt)

	case reflect.Slice:
		t.Type = Array
		it := Type{}.resolve(c, rt.Elem())
		t.Items = &it

	default:
		panic("cannot resolve schema for type: " + rt.Kind().String())
	}
	return t
}
