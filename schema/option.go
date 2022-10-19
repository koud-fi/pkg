package schema

import "reflect"

type Option func(*config)

func CustomType(packagePath, typeName string, fn func(reflect.Type) Type) Option {
	return func(c *config) { c.customTypes[typeKey{packagePath, typeName}] = fn }
}

func Tags(tag ...string) Option { return func(c *config) { c.tags = tag } }

type config struct {
	customTypes map[typeKey]func(reflect.Type) Type
	tags        []string
}

type typeKey struct {
	packagePath string
	typeName    string
}
