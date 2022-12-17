package pk

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/koud-fi/pkg/blob"

	"golang.org/x/crypto/sha3"
)

const uidByteCount = 15

var uidSrcSeparator = []byte{0}

type UID struct {
	key [uidByteCount]byte
}

type Salt []byte

func NewUID(src ...any) (UID, error) {
	var u UID
	if len(src) == 0 {
		if _, err := rand.Read(u.key[:]); err != nil {
			return u, err
		}
	} else {
		h := sha3.NewShake128()
		for i := range src {
			v := src[i]
			if salt, ok := v.(Salt); ok {
				v = []byte(salt)
			} else if i > 0 {
				h.Write(uidSrcSeparator)
			}
			switch v := v.(type) {
			case []byte:
				h.Write(v)

			case blob.Blob:
				if b, err := v.Open(); err != nil {
					return u, err
				} else {
					if _, err := io.Copy(h, b); err != nil {
						return u, err
					}
				}
			default:
				fmt.Fprint(h, v)
			}
		}
		h.Read(u.key[:])
	}
	return u, nil
}

func ParseUID(s string) (UID, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return UID{}, fmt.Errorf("malformed UID, %w", err)
	}
	return ParseUIDBytes(b)
}

func ParseUIDBytes(b []byte) (UID, error) {
	if len(b) != uidByteCount {
		return UID{}, errors.New("malformed UID, unexpected length")
	}
	var t UID
	copy(t.key[:], b)
	return t, nil
}

func (u UID) Bytes() []byte { return u.key[:] }

func (u UID) String() string { return base64.URLEncoding.EncodeToString(u.key[:]) }

func (t UID) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, uidByteCount+2))
	buf.WriteByte('"')
	buf.Write(t.key[:])
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func (u *UID) UnmarshalJSON(data []byte) (err error) {
	n := len(data)
	if n >= 2 && data[0] == '"' && data[n-1] == '"' {
		data = data[1 : n-1]
	}
	*u, err = ParseUIDBytes(data)
	return
}
