package blob

import "strings"

const (
	Default Domain = ""

	refDomainSeparator = ":"
)

var ZeroRef = NewRef(Default)

type (
	Domain string
	Ref    []string
)

func NewRef(d Domain, ref ...string) Ref {
	refs := make([]string, 0, len(ref)+1)
	for _, ref := range ref {
		refs = append(refs, strings.Split(ref, refDomainSeparator)...)
	}
	if d != Default {
		if strings.Contains(string(d), refDomainSeparator) {
			panic(`blob.NewRef: invalid domain, contains ":"`)
		}
	}
	return Ref(refs)
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
