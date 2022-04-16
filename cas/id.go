package cas

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"

	"golang.org/x/crypto/sha3"
)

const (
	idLen       = 20
	idHexLen    = 30
	idByteCount = 15
)

var (
	// ErrInvalidID is used when ID parsing fails for some reason.
	ErrInvalidID = errors.New("invalid ID")

	defaultIDSalt = []byte{
		11, 134, 201, 40, 86, 208, 163, 24,
		31, 153, 181, 130, 34, 104, 176, 24,
		89, 26, 146, 221, 123, 106, 122, 55,
		248, 125, 115, 77, 24, 63, 203, 135,
		191, 181, 188, 81, 186, 16, 158, 188,
		2, 18, 255, 200, 5, 211, 56, 238,
		47, 109, 15, 180, 46, 31, 67, 236,
		138, 250, 175, 149, 50, 155, 94, 97,
	}
	idValidator = regexp.MustCompile(fmt.Sprintf("^[a-zA-Z0-9_-]{%d}$", idLen))
)

// ID is a 120-bit (15 byte) "content ID" for arbitrary binary data, with
// collision chance of 1 in ~2.66 trillion on trillion unique items.
type ID string

// NewID reads all data from r and calculates a ID for it.
func NewID(r io.Reader) (ID, error) {
	h := sha3.NewShake128()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	h.Write(defaultIDSalt)

	var buf [idByteCount]byte
	h.Read(buf[:])
	return ID(base64.URLEncoding.EncodeToString(buf[:])), nil
}

// NewIDFromBytes calculates a ID for the given b.
func NewIDFromBytes(b []byte) ID {
	id, _ := NewID(bytes.NewReader(b))
	return id
}

// ParseID parses s into a valid ID, s can be either standard or hex encoded ID.
func ParseID(s string) (ID, error) {
	buf := make([]byte, idByteCount)
	switch len(s) {
	case idLen:
		base64.URLEncoding.Decode(buf, []byte(s))
	case idHexLen:
		hex.Decode(buf, []byte(s))
	default:
		return "", ErrInvalidID
	}
	return ID(base64.URLEncoding.EncodeToString(buf)), nil
}

// Hex returns a hex encoded presentation of a ID.
// This should be used when using IDs as filenames on case-insensitive filesystems.
func (id ID) Hex() string {
	buf, err := base64.URLEncoding.DecodeString(string(id))
	if err != nil {
		panic("malformed ID")
	}
	return hex.EncodeToString(buf)
}

func (id ID) String() string { return string(id) }
