package rpcapi

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/koud-fi/pkg/assign"
	"github.com/koud-fi/pkg/errx"
)

type Arguments interface {
	Get(key string) any
}

type ArgumentMap map[string]any

func (m ArgumentMap) Get(key string) any { return m[key] }

type URLValueArguments url.Values

func (u URLValueArguments) Get(key string) any {
	vs, ok := u[key]
	if !ok || len(vs) == 0 {
		return nil
	}
	if len(vs) == 1 {
		return vs[0]
	}
	return vs
}

// ApplyArguments populates dst (a pointer to a struct) with values from args, using the converter.
// It uses both the original field name and a normalized version (e.g. lower-case)
// so that input keys can be in mixed case.
func ApplyArguments(dst any, converter *assign.Converter, args Arguments) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errx.New("dst must be a pointer to a struct")
	}
	rv = rv.Elem()
	rt := rv.Type()
	for i := range rt.NumField() {
		f := rt.Field(i)

		// Handle embedded struct fields
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			embedded := reflect.New(f.Type).Interface()
			if err := ApplyArguments(embedded, converter, args); err != nil {
				return errx.Fmt("embedded field %q: %w", f.Name, err)
			}
			rv.Field(i).Set(reflect.ValueOf(embedded).Elem())
			continue
		}
		var (
			key = f.Name
			val = args.Get(key)
		)
		if val == nil {
			normalized := normalizeArgumentKey(key)
			if normalized != key {
				val = args.Get(normalized)
			}
		}
		if val == nil {
			continue
		}
		conv, err := converter.Convert(val, f.Type)
		if err != nil {
			return errx.Fmt("field %q: %w", key, err)
		}
		rv.Field(i).Set(conv)
	}
	return nil
}

// normalizeArgumentKey returns a canonical form of the key (here simply lower-case).
func normalizeArgumentKey(s string) string {
	return strings.ToLower(s)
}

type CombinedArguments []Arguments

func (c CombinedArguments) Get(key string) any {
	for _, args := range c {
		if v := args.Get(key); v != nil {
			return v
		}
	}
	return nil
}
