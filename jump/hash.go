package jump

import "io"

// Hash consistently chooses a hash bucket number in the range [0, numBuckets) for the given key
// numBuckets must be >= 1
func Hash(key uint64, numBuckets int32) int32 {
	var (
		b int64 = -1
		j int64
	)
	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}
	return int32(b)
}

// HashString works like Hash but takes string as a key,
// using KeyHasher to make it compatible with Hash
func HashString(key string, numBuckets int32, h KeyHasher) int32 {
	h.Reset()
	if _, err := io.WriteString(h, key); err != nil {
		panic(err)
	}
	return Hash(h.Sum64(), numBuckets)
}

// KeyHasher is a subset of hash.Hash64
type KeyHasher interface {
	io.Writer
	Reset()
	Sum64() uint64
}
