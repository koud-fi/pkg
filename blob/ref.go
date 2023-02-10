package blob

import "strings"

const (
	Default Domain = ""

	refDomainSeparator = ":"
)

type (
	Domain string
	Ref    []string
)

func NewRef(s string) Ref {
	return Ref(strings.Split(s, refDomainSeparator))
}

func (r Ref) Domain() Domain {
	if len(r) > 1 {
		return Domain(r[0])
	}
	return ""
}

func (r Ref) Ref() Ref {
	if len(r) > 1 {
		return Ref(r[1:])
	}
	return Ref{}
}

func (r Ref) String() string { return strings.Join(r, refDomainSeparator) }

func (r Ref) Bytes() []byte { return []byte(r.String()) }

// TODO: MarshalJSON/UnmarshalJSON
