package pk

import (
	"fmt"
	"sync"
	"time"
)

type TIDSource struct {
	mu sync.Mutex
	ts int64 // Current "time step"
	n  uint  // User defined number of the source
	c  uint  // Counter in the current "time step"
}

func NewTIDSource(now time.Time, n uint) *TIDSource {
	if n == 0 || n >= tidNumberMax {
		panic(fmt.Sprintf("n must be greater than 0 and less than %d", tidNumberMax))
	}
	return &TIDSource{
		ts: nowTimeStep(now),
		n:  n,
		c:  tidCounterMax, // Don't generate any new IDs for the current time step
	}
}

func (src *TIDSource) Next(now time.Time) TID        { return src.next(now, false) }
func (src *TIDSource) NextVirtual(now time.Time) TID { return src.next(now, true) }

func (src *TIDSource) next(now time.Time, virtual bool) TID {
	src.mu.Lock()
	defer src.mu.Unlock()

	nowTs := nowTimeStep(now)
	if nowTs <= src.ts {
		if src.c == tidCounterMax {
			src.ts++
			time.Sleep(time.Duration((src.ts - nowTs) * tidTimeStep))
			src.c = 1
		} else {
			src.c++
		}
	} else {
		src.ts = nowTs
		src.c = 1
	}
	return NewTID(time.Unix(src.ts, 0), src.n, src.c-1, virtual)
}

func nowTimeStep(now time.Time) int64 {
	return now.UTC().UnixNano() / tidTimeStep
}
