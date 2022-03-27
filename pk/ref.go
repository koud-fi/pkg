package pk

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	refKeySeparator    = "/"
	refParamsSeparator = ":"
)

type Ref struct {
	scheme Scheme
	params string
	key    string
}

type Scheme string

func NewRef(s Scheme, params, key string) (Ref, error) {
	r := Ref{scheme: s, params: params, key: strings.TrimPrefix(key, refKeySeparator)}
	if err := r.Validate(); err != nil {
		return Ref{}, err
	}
	return r, nil
}

func ParseRef(s string) (ref Ref, err error) {
	err = ref.parse(s)
	return
}

func (r Ref) Scheme() Scheme { return r.scheme }
func (r Ref) Params() string { return r.params }
func (r Ref) Key() string    { return r.key }

func (r Ref) Validate() error {

	// TODO: improve validation

	if len(r.scheme) == 0 {
		return fmt.Errorf("ref: %w; empty schema", os.ErrInvalid)
	}
	return nil
}

func (r Ref) String() string { return r.buf().String() }
func (r Ref) Bytes() []byte  { return r.buf().Bytes() }

func (r Ref) buf() *bytes.Buffer {
	if err := r.Validate(); err != nil {
		return nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, r.bufSize()))
	r.encode(buf)
	return buf
}

func (r Ref) bufSize() (size int) {
	size = len(r.scheme) + len(r.key) + 1
	if r.params != "" {
		size += len(r.params) + 1
	}
	return
}

func (r Ref) encode(buf *bytes.Buffer) {
	buf.WriteString(refKeySeparator + string(r.scheme))
	if r.params != "" {
		buf.WriteString(refParamsSeparator)
		buf.WriteString(r.params)
	}
	buf.WriteString(refKeySeparator)
	buf.WriteString(escapeKey(r.key))
}

func (r *Ref) parse(s string) (err error) {
	s = strings.TrimPrefix(s, refKeySeparator)
	if s == "" {
		return fmt.Errorf("ref: %w; empty", os.ErrInvalid)
	}
	parts := strings.SplitN(s, refKeySeparator, 2)
	if len(parts) != 2 {
		return fmt.Errorf("ref: %w; (%s) missing separator", os.ErrInvalid, s)
	}
	unescapedKey, err := unescapeKey(parts[1])
	if err != nil {
		return fmt.Errorf("ref: unescape failed; %w", err)
	}
	var (
		paramsParts = strings.SplitN(parts[0], refParamsSeparator, 2)
		scheme      = Scheme(paramsParts[0])
	)
	switch len(paramsParts) {
	case 1:
		*r, err = NewRef(scheme, "", refKeySeparator)
	case 2:
		*r, err = NewRef(scheme, paramsParts[1], unescapedKey)
	}
	return
}

func escapeKey(key string) string            { return key }
func unescapeKey(key string) (string, error) { return key, nil }

func (r Ref) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(make([]byte, 0, r.bufSize()+2))
	buf.WriteRune('"')
	r.encode(buf)
	buf.WriteRune('"')
	return buf.Bytes(), nil
}

func (r *Ref) UnmarshalJSON(data []byte) error {
	return r.parse(string(bytes.Trim(data, `"`)))
}
