package schema

import "reflect"

type Option func(*config)

func WithCustomType(pkgPath, name string, fn func(reflect.Type) Type) Option {
	return func(c *config) { c.customTypes[typeKey{pkgPath, name}] = fn }
}

// WithTags allows overriding which struct tags to honor (e.g. "json", "xml").
func WithTags(tags ...string) Option {
	return func(c *config) { c.tags = tags }
}

type config struct {
	customTypes map[typeKey]func(reflect.Type) Type
	tags        []string
	inProgress  map[reflect.Type]*Type
}

type typeKey struct {
	pkgPath string
	name    string
}
