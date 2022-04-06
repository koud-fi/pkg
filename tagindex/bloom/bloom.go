package bloom

import (
	"encoding/binary"
	"fmt"

	"github.com/cespare/xxhash"
)

type Filter [4]uint64

func New(data []uint64, k int) Filter {
	var (
		h = xxhash.New()
		b [8]byte
		f Filter
	)
	for i, n := range data {
		if i > 0 {
			h.Reset()
		}
		for j := 0; j < k; j++ {
			binary.LittleEndian.PutUint64(b[:], n)
			h.Write(b[:])
			var (
				a = h.Sum64() % 256
				b = a / 64
			)
			f[b] = f[b] | 1<<(a-b*64)
		}
	}
	return f
}

func New32(data []uint32, k int) Filter {
	ns := make([]uint64, len(data))
	for i := range data {
		ns[i] = uint64(data[i])
	}
	return New(ns, k)
}

func (t Filter) Contains(o Filter) bool {
	return t[0]&o[0] == o[0] && t[1]&o[1] == o[1] && t[2]&o[2] == o[2] && t[3]&o[3] == o[3]
}

func (t Filter) String() string {
	return fmt.Sprintf("%064b %064b %064b %064b", t[0], t[1], t[2], t[3])
}
