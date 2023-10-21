package file

import (
	"crypto"
	"encoding/hex"

	"github.com/koud-fi/pkg/blob"
)

// TODO: support custom hash functions

func Digests(h ...crypto.Hash) Option {
	return func(a *Attributes, b blob.Blob, contentType string) error {
		if a.IsDir {
			return nil
		}

		// TODO: avoid full memory copy

		var data []byte
		for _, hashType := range h {
			if !hashType.Available() {
				continue
			}
			if data == nil {
				var err error
				if data, err = blob.Bytes(b); err != nil {
					return err
				}
			}
			hash := hashType.New()
			if _, err := hash.Write(data); err != nil {
				return err
			}
			if a.Digest == nil {
				a.Digest = make(map[string]string, len(h))
			}
			a.Digest[hashType.String()] = hex.EncodeToString(hash.Sum(nil))
		}
		return nil
	}
}
