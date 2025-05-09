package pk

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

const tidTimeStep = int64(time.Millisecond)

// TID is an 63-bit "temporal" ID that is generated from current time and a magic number.
// 1 million new TIDs can be generated each millisecond, they will run out around year 2262.
type TID struct {
	value int64 // We use a signed int64 as some datastores don't support unsigned integers.
}

// New creates a TID from given time and number, N must be a positive integer under 1 million.
func NewTID(now time.Time, n uint) TID {
	return newTID(now.UTC().UnixNano()/tidTimeStep, int64(n))
}

func newTID(ts, n int64) TID {
	if n < 0 || n >= tidTimeStep {
		panic("tid: invalid 'n' for TID")
	}
	return TID{value: ts*tidTimeStep + n}
}

// Set sets the TID to a raw int64 value, should be used with care.
// This exists to enable the implementation of custom marshalling.
func (t *TID) Set(value int64) { t.value = value }

// ParseTID converts string generated by the TIDs String function into a TID.
func ParseTID(s string) (TID, error) {
	return ParseTIDBytes([]byte(s))
}

// ParseTIDBytes converts raw bytes into a TID.
func ParseTIDBytes(raw []byte) (TID, error) {
	var b [8]byte
	n, err := base64.RawURLEncoding.Decode(b[:], bytes.Trim(raw, `"`))
	if err != nil {
		return TID{}, fmt.Errorf("tid: decode: %w", err)
	}
	switch {
	case n > 8:
		return TID{}, fmt.Errorf("tid: invalid length")
	case n < 8:
		// Shift bytes to end and pad with leading zeros
		copy(b[8-n:], b[:n])
		for i := range 8 - n {
			b[i] = 0
		}
	}
	return TID{value: int64(binary.BigEndian.Uint64(b[:]))}, nil
}

// IsZero returns true if the TID is zero (uninitialized).
func (t TID) IsZero() bool { return t.value == 0 }

// Time converts TID into a unique time.Time with millisecond precision,
// the micro/nanosecond part is effectively noise.
func (t TID) Time() time.Time {
	s := int64(t.value) / int64(time.Second)
	n := int64(t.value) % int64(time.Second)
	return time.Unix(s, n)
}

// Value returns the TID as an int64.
func (t TID) Value() int64 { return t.value }

// String converts TID to a base-64 string.
func (t TID) String() string {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(t.value))

	// Find first non-zero byte to avoid encoding leading zeros
	i := 0
	for i < len(b) && b[i] == 0 {
		i++
	}
	return base64.RawURLEncoding.EncodeToString(b[i:])
}

// MarshalJSON implements the json.Marshaler interface.
func (t TID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TID) UnmarshalJSON(b []byte) error {
	tid, err := ParseTIDBytes(b)
	if err != nil {
		return fmt.Errorf("tid: unmarshal: %w", err)
	}
	*t = tid
	return nil
}

type TIDSource struct {
	mu sync.Mutex
	ts int64
	n  int

	// TODO: support for node ID
}

func NewTIDSource(now time.Time) *TIDSource {
	return &TIDSource{
		ts: now.UTC().UnixNano() / tidTimeStep,
		n:  1000,
	}
}

func (t *TIDSource) Next(now time.Time) TID {
	t.mu.Lock()
	defer t.mu.Unlock()

	nowTs := now.UTC().UnixNano() / int64(time.Millisecond)
	if nowTs <= t.ts {
		if t.n == 1000 {
			t.ts++
			time.Sleep(time.Duration((t.ts - nowTs) * tidTimeStep))
			t.n = 1
		} else {
			t.n++
		}
	} else {
		t.ts = nowTs
		t.n = 1
	}
	return newTID(t.ts, int64(t.n-1)*1000)
}
