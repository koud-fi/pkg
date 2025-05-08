package assign

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/koud-fi/pkg/errx"
)

var ErrUnsupportedConversion = errors.New("unsupported conversion")

type ConverterFunc func(in any, target reflect.Type) (reflect.Value, error)

type Converter struct {
	funcs []ConverterFunc
}

func NewDefaultConverter() *Converter {
	conv := new(Converter)
	conv.Register(ConvertJSON) // Do this first to handle json.Unmarshaler before primitive conversions.
	conv.Register(ConvertPrimitive)
	conv.Register(ConvertSlice(conv))
	conv.Register(ConvertMap(conv))
	conv.Register(ConvertStruct(conv))
	return conv
}

func (c *Converter) Register(conv ConverterFunc) {
	c.funcs = append(c.funcs, conv)
}

func (c *Converter) Convert(in any, target reflect.Type) (reflect.Value, error) {
	for _, fn := range c.funcs {
		v, err := fn(in, target)
		if err == nil {
			return v, nil
		}
		if !errors.Is(err, ErrUnsupportedConversion) {
			return reflect.Value{}, errx.E(err)
		}
	}
	return reflect.Value{}, errx.Fmt("can't convert %v to %v", reflect.TypeOf(in), target)
}

func ConvertPrimitive(in any, target reflect.Type) (reflect.Value, error) {
	v := reflect.ValueOf(in)
	if v.Type().AssignableTo(target) {
		return v, nil
	}
	if v.Type().ConvertibleTo(target) {
		return v.Convert(target), nil
	}
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := fmt.Sprint(in)
		n, err := strconv.ParseInt(s, 10, target.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("can't parse %q as int: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := fmt.Sprint(in)
		n, err := strconv.ParseUint(s, 10, target.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("can't parse %q as uint: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Float32, reflect.Float64:
		s := fmt.Sprint(in)
		n, err := strconv.ParseFloat(s, target.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("can't parse %q as float: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Bool:
		s := fmt.Sprint(in)
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("can't parse %q as bool: %w", s, err)
		}
		return reflect.ValueOf(b), nil

	case reflect.String:
		return reflect.ValueOf(fmt.Sprint(in)), nil

	case reflect.Pointer:
		elType := target.Elem()
		elVal, err := ConvertPrimitive(in, elType)
		if err != nil {
			return reflect.Value{}, errx.E(err)
		}
		ptr := reflect.New(elType)
		ptr.Elem().Set(elVal)
		return ptr, nil
	}
	return reflect.Value{}, errx.E(ErrUnsupportedConversion)
}

func ConvertSlice(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Slice {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		rv := reflect.ValueOf(in)
		if rv.Kind() != reflect.Slice {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		n := rv.Len()
		out := reflect.MakeSlice(target, n, n)
		for i := range n {
			elem, err := conv.Convert(rv.Index(i).Interface(), target.Elem())
			if err != nil {
				return reflect.Value{}, errx.Fmt("slice index %d: %w", i, err)
			}
			out.Index(i).Set(elem)
		}
		return out, nil
	}
}

func ConvertMap(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Map {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		rv := reflect.ValueOf(in)
		if rv.Kind() != reflect.Map {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		out := reflect.MakeMap(target)
		for _, key := range rv.MapKeys() {
			newKey, err := conv.Convert(key.Interface(), target.Key())
			if err != nil {
				return reflect.Value{}, errx.Fmt("map key %v: %w", key, err)
			}
			newVal, err := conv.Convert(rv.MapIndex(key).Interface(), target.Elem())
			if err != nil {
				return reflect.Value{}, errx.Fmt("map value for key %v: %w", key, err)
			}
			out.SetMapIndex(newKey, newVal)
		}
		return out, nil
	}
}

func ConvertStruct(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Struct {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		m, ok := in.(map[string]any)
		if !ok {
			return reflect.Value{}, errx.E(ErrUnsupportedConversion)
		}
		out := reflect.New(target).Elem()
		for i := range target.NumField() {
			f := target.Field(i)

			switch {
			case f.Anonymous && f.Type.Kind() == reflect.Struct:
				embeddedVal, err := conv.Convert(m, f.Type)
				if err != nil {
					return reflect.Value{}, errx.Fmt("embedded field %q: %w", f.Name, err)
				}
				out.Field(i).Set(embeddedVal)
				continue

			default:
				if val, exists := m[f.Name]; exists {
					v, err := conv.Convert(val, f.Type)
					if err != nil {
						return reflect.Value{}, errx.Fmt("field %q: %w", f.Name, err)
					}
					out.Field(i).Set(v)
				}
			}
		}
		return out, nil
	}
}

func ConvertJSON(in any, target reflect.Type) (reflect.Value, error) {
	var data []byte
	switch v := in.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	case json.RawMessage:
		data = []byte(v)
	default:
		return reflect.Value{}, errx.E(ErrUnsupportedConversion)
	}
	if Implements[json.Unmarshaler](target) {
		ptr := reflect.New(target)
		if err := ptr.Interface().(json.Unmarshaler).UnmarshalJSON(data); err != nil {
			return reflect.Value{}, fmt.Errorf("json unmarshal error: %w", err)
		}
		return ptr.Elem(), nil
	}
	switch target.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		ptr := reflect.New(target)
		if err := json.Unmarshal(data, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("json unmarshal error: %w", err)
		}
		return ptr.Elem(), nil
	}
	return reflect.Value{}, errx.E(ErrUnsupportedConversion)
}
