package assign

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	DefaultConverter = NewDefaultConverter()

	ErrUnsupportedConversion = errors.New("unsupported conversion")
)

type ConverterFunc func(in any, target reflect.Type) (reflect.Value, error)

type Converter struct {
	funcs []ConverterFunc
}

func NewDefaultConverter() *Converter {
	conv := new(Converter)
	return conv.Register(ConvertPrimitive).
		Register(ConvertSlice(conv)).
		Register(ConvertMap(conv)).
		Register(ConvertStruct(conv)).
		Register(ConvertJSON)
}

func (c Converter) Register(conv ConverterFunc) *Converter {
	c.funcs = append(c.funcs, conv)
	return &c
}

func (c *Converter) Convert(in any, target reflect.Type) (reflect.Value, error) {
	for _, fn := range c.funcs {
		v, err := fn(in, target)
		if err == nil {
			return v, nil
		}
		if !errors.Is(err, ErrUnsupportedConversion) {
			return reflect.Value{}, err
		}
	}
	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", reflect.TypeOf(in), target)
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
			return reflect.Value{}, fmt.Errorf("cannot parse %q as int: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := fmt.Sprint(in)
		n, err := strconv.ParseUint(s, 10, target.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot parse %q as uint: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Float32, reflect.Float64:
		s := fmt.Sprint(in)
		n, err := strconv.ParseFloat(s, target.Bits())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot parse %q as float: %w", s, err)
		}
		return reflect.ValueOf(n).Convert(target), nil

	case reflect.Bool:
		s := fmt.Sprint(in)
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot parse %q as bool: %w", s, err)
		}
		return reflect.ValueOf(b), nil

	case reflect.String:
		return reflect.ValueOf(fmt.Sprint(in)), nil
	}
	return reflect.Value{}, ErrUnsupportedConversion
}

func ConvertSlice(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Slice {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		rv := reflect.ValueOf(in)
		if rv.Kind() != reflect.Slice {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		n := rv.Len()
		out := reflect.MakeSlice(target, n, n)
		for i := 0; i < n; i++ {
			elem, err := conv.Convert(rv.Index(i).Interface(), target.Elem())
			if err != nil {
				return reflect.Value{}, fmt.Errorf("slice index %d: %w", i, err)
			}
			out.Index(i).Set(elem)
		}
		return out, nil
	}
}

func ConvertMap(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Map {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		rv := reflect.ValueOf(in)
		if rv.Kind() != reflect.Map {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		out := reflect.MakeMap(target)
		for _, key := range rv.MapKeys() {
			newKey, err := conv.Convert(key.Interface(), target.Key())
			if err != nil {
				return reflect.Value{}, fmt.Errorf("map key %v: %w", key, err)
			}
			newVal, err := conv.Convert(rv.MapIndex(key).Interface(), target.Elem())
			if err != nil {
				return reflect.Value{}, fmt.Errorf("map value for key %v: %w", key, err)
			}
			out.SetMapIndex(newKey, newVal)
		}
		return out, nil
	}
}

func ConvertStruct(conv *Converter) ConverterFunc {
	return func(in any, target reflect.Type) (reflect.Value, error) {
		if target.Kind() != reflect.Struct {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		m, ok := in.(map[string]any)
		if !ok {
			return reflect.Value{}, ErrUnsupportedConversion
		}
		out := reflect.New(target).Elem()
		for i := 0; i < target.NumField(); i++ {
			f := target.Field(i)
			if val, exists := m[f.Name]; exists {
				conv, err := conv.Convert(val, f.Type)
				if err != nil {
					return reflect.Value{}, fmt.Errorf("field %q: %w", f.Name, err)
				}
				out.Field(i).Set(conv)
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
		return reflect.Value{}, ErrUnsupportedConversion
	}
	switch target.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		ptr := reflect.New(target)
		if err := json.Unmarshal(data, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("json unmarshal error: %w", err)
		}
		return ptr.Elem(), nil
	}
	return reflect.Value{}, ErrUnsupportedConversion
}
