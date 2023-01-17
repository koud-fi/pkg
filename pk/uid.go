package pk

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/koud-fi/pkg/blob"

	"golang.org/x/crypto/sha3"
)

const (
	uidRawLen   = 15
	idBase64Len = 20
	idHexLen    = 30
)

var uidSrcSeparator = []byte{0}

// UID is a 120-bit (15 byte) "content ID" for arbitrary binary data, with
// collision chance of 1 in ~2.66 trillion on trillion unique items.
type UID struct{ key [uidRawLen]byte }

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
	var (
		buf = make([]byte, uidRawLen)
		err error
	)
	switch len(s) {
	case idBase64Len:
		_, err = base64.URLEncoding.Decode(buf, []byte(s))
	case idHexLen:
		_, err = hex.Decode(buf, []byte(s))
	default:
		err = errors.New("unknown encoding")
	}
	if err != nil {
		return UID{}, fmt.Errorf("malformed UID, %w", err)
	}
	return ParseUIDBytes(buf)
}

func ParseUIDBytes(b []byte) (UID, error) {
	if len(b) != uidRawLen {
		return UID{}, errors.New("malformed UID, unexpected length")
	}
	var t UID
	copy(t.key[:], b)
	return t, nil
}

func (u UID) Bytes() []byte { return u.key[:] }

// Hex returns a hex encoded presentation of a ID.
// This should be used when using IDs as filenames on case-insensitive filesystems.
func (u UID) Hex() string { return hex.EncodeToString(u.key[:]) }

func (u UID) String() string { return base64.URLEncoding.EncodeToString(u.key[:]) }

func (t UID) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, uidRawLen+2))
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
