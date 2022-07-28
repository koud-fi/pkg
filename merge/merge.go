package merge

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Values(dst any, valueFn func(key string) []string) error {
	return applyValues(reflect.ValueOf(dst), valueFn)
}

func applyValues(dst reflect.Value, valueFn func(key string) []string) error {
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	var (
		vType = dst.Type()
		vLen  = dst.NumField()
		err   error
	)
	for i := 0; i < vLen && err == nil; i++ {
		var (
			f     = vType.Field(i)
			fVal  = dst.Field(i)
			fName = resolveFieldName(f)
		)
		if f.Anonymous {
			applyValues(fVal, valueFn)
			continue
		}
		vs := valueFn(fName)
		if len(vs) == 0 {
			continue
		}
		if fVal.CanAddr() && fVal.CanInterface() {
			vi := fVal.Addr().Interface()
			if p, ok := vi.(interface{ Parse(string) error }); ok {
				err = p.Parse(vs[0])
				continue
			}
		}
		fKind := f.Type.Kind()
		switch fKind {
		case reflect.String:
			fVal.SetString(vs[0])

		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			var i int64
			if i, err = strconv.ParseInt(vs[0], 10, f.Type.Bits()); err == nil {
				fVal.SetInt(i)
			}
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			var u uint64
			if u, err = strconv.ParseUint(vs[0], 10, f.Type.Bits()); err == nil {
				fVal.SetUint(u)
			}
		case reflect.Float64, reflect.Float32:
			var fl float64
			if fl, err = strconv.ParseFloat(vs[0], f.Type.Bits()); err == nil {
				fVal.SetFloat(fl)
			}
		case reflect.Bool:
			var b bool
			if b, err = strconv.ParseBool(vs[0]); err == nil {
				fVal.SetBool(b)
			}

		// TODO: slice handling

		// TODO: struct handling

		default:
			return fmt.Errorf("unsupported kind %v for value: %s", fKind, fName)
		}
	}
	return err
}

func resolveFieldName(f reflect.StructField) string {
	name := f.Name
	if tag, ok := f.Tag.Lookup("json"); ok {
		parts := strings.SplitN(tag, ",", 2)
		if parts[0] != "" {
			name = parts[0]
		}
	}
	return name
}
