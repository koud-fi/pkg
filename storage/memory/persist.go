package memory

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/blob/localfile"
	"github.com/koud-fi/pkg/rx"
)

func Save(path string, s *Storage, perm os.FileMode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: avoid full memory copy by streaming the data

	buf := bytes.NewBuffer(nil)
	if err := rx.ForEachN(rx.SliceIter(s.data...), func(p rx.Pair[string, []byte], i int) error {
		if bytes.ContainsRune(p.Value, '\n') {
			return errors.New("values must not contain linebreaks")
		}
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(url.PathEscape(p.Key))
		buf.WriteByte(' ')
		buf.Write(p.Value)
		return nil

	}); err != nil {
		return err
	}
	return localfile.Write(path, blob.FromBytes(buf.Bytes()), perm)
}

func Load(path string, s *Storage) error {

	// TODO: use storage internals to make data writing faster

	ctx := context.Background()
	return rx.ForEachN(blob.LineIter(localfile.New(path)), func(data string, line int) error {
		splitAt := strings.IndexByte(data, ' ')
		if splitAt == -1 {
			return fmt.Errorf("malformed data at line %d", line)
		}
		key, err := url.PathUnescape(data[splitAt:])
		if err != nil {
			return fmt.Errorf("malformed key at line %d: %w", line, err)
		}
		s.Set(ctx, key, strings.NewReader(data[splitAt+1:]))
		return nil
	})
}
