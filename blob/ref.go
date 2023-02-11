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

func NewRef(d Domain, ref string) Ref {
	if d == Default {
		return ParseRef(ref)
	}
	return append(Ref{string(d)}, ParseRef(ref)...)
}

func ParseRef(ref string) Ref {
	return Ref(strings.Split(ref, refDomainSeparator))
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
