package blob

import "strings"

const DomainSeparator = ":"

type Ref []string

func ParseRef(s string) Ref {
	return Ref(strings.Split(s, DomainSeparator))
}

func (r Ref) Domain() string {
	if len(r) > 1 {
		return r[0]
	}
	return ""
}

func (r Ref) Ref() Ref {
	if len(r) > 1 {
		return Ref(r[1:])
	}
	return Ref{}
}

func (r Ref) String() string { return strings.Join(r, DomainSeparator) }

// TODO: MarshalJSON/UnmarshalJSON
