package pk

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/sha3"
)

const uidByteCount = 15

var uidSrcSeparator = []byte{0}

type UID struct {
	key [uidByteCount]byte
}

type Salt []byte

func NewUID(src ...any) (t UID) {
	if len(src) == 0 {
		if _, err := rand.Read(t.key[:]); err != nil {
			panic("pk.NewUID: " + err.Error())
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
			switch v := src[i].(type) {
			case []byte:
				h.Write(v)
			default:
				fmt.Fprint(h, v)
			}
		}
		h.Read(t.key[:])
	}
	return
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

func (t UID) Bytes() []byte { return t.key[:] }

func (t UID) String() string { return base64.URLEncoding.EncodeToString(t.key[:]) }
