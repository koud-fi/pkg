package pk

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const tidTimeStep = int64(time.Millisecond)

// TID is an 63-bit "temporal" ID that is generated from current time and a magic number.
// 1 million new TIDs can be generated each millisecond, they will run out around year 2262.
type TID int64

// New creates a TID from given time and number, N must be a positive integer under 1 million.
func NewTID(now time.Time, n uint) TID {
	return newTID(now.UTC().UnixNano()/tidTimeStep, int64(n))
}

func newTID(ts, n int64) TID {
	if n < 0 || n >= tidTimeStep {
		panic("pk: invalid 'n' for TID")
	}
	return TID(ts*tidTimeStep + n)
}

// Parse convert string generated by the TIDs String function into a TID.
func ParseTID(s string) (TID, error) {
	n, err := strconv.ParseInt(s, 36, 64)
	if err != nil {
		return 0, fmt.Errorf("pk: %w", err)
	}
	return TID(n), nil
}

// Time converts TID into a unique time.Time with millisecond precision,
// the micro/nanosecond part is effectively noise.
func (t TID) Time() time.Time {
	s := int64(t) / int64(time.Second)
	n := int64(t) % int64(time.Second)
	return time.Unix(s, n)
}

// Int64 is just a cast from TID to int64.
func (t TID) Int64() int64 { return int64(t) }

// String converts TID to a base-36 string.
func (t TID) String() string { return strconv.FormatInt(int64(t), 36) }

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
