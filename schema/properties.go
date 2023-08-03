package schema

import "reflect"

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
