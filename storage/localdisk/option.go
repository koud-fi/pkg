package localdisk

import (
	"crypto"
	"os"
)

type Option func(*Storage)

func Buckets(levels ...int) Option {
	return func(s *Storage) {
		s.bucketLevels = levels
		s.bucketPrefixLen = 0
		for _, l := range levels {
			s.bucketPrefixLen += l + 1
		}
	}
}

func BucketHash(h crypto.Hash) Option { return func(s *Storage) { s.bucketHash = &h } }
func DirPerm(m os.FileMode) Option    { return func(s *Storage) { s.dirPerm = m } }
func FilePerm(m os.FileMode) Option   { return func(s *Storage) { s.filePerm = m } }
