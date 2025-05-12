package feistel_test

import (
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/koud-fi/pkg/secret/feistel"
)

func TestFeistel(t *testing.T) {
	keys := feistel.KeysFromString("salakala")
	for i := range 100 {
		var (
			original  = uint64(i)
			encrypted = feistel.Encrypt(original, keys[:])
			decrypted = feistel.Decrypt(encrypted, keys[:])
			//pretty    = strings.ToUpper(strconv.FormatInt(int64(encrypted), 36))
		)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, encrypted)
		pretty := base64.RawURLEncoding.EncodeToString(buf)

		t.Logf("%20d -> %20d | %20s -> %20d\n", original, encrypted, pretty, decrypted)

		if original != decrypted {
			t.Fatal("Decrypted value does not match the original.")
		}
	}
}
