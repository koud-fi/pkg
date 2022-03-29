package protoserver

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/cache/blobcache"
	"github.com/koud-fi/pkg/pk"
	"github.com/koud-fi/pkg/schema"
)

const FileScheme = "file"

func RegisterFile(cacheStorage blob.Storage) {
	c := blobcache.New(cacheStorage)

	Register(FileScheme,
		TransformFetcher(func(ctx context.Context, key string, ref pk.Ref) (any, error) {
			cacheKey := ref.String()
			if key == "" {
				key = schema.MasterFileKey
			} else if key != schema.MasterFileKey {
				cacheKey += "." + key
			}
			return c.Resolve(ctx, cacheKey, blob.Func(func() (io.ReadCloser, error) {
				v, err := Fetch(ctx, ref)
				if err != nil {
					return nil, err
				}
				fileNode, ok := v.(schema.FileNode)
				if !ok {
					return nil, fmt.Errorf("%w: is not a file node", fs.ErrNotExist)
				}
				files, err := fileNode.Files()
				if err != nil {
					return nil, fmt.Errorf("error resolving files: %w", err)
				}
				if key == "" {
					key = schema.MasterFileKey
				}
				f, ok := files[key]
				if !ok {
					return nil, fmt.Errorf("%w: %s", fs.ErrNotExist, key)
				}
				rc, err := f.Open()
				if err != nil {
					return nil, fmt.Errorf("error opening file: %w", err)
				}
				if f, ok := rc.(*os.File); ok {
					return f, blobcache.NoCache
				}
				return rc, nil
			})).Open()
		}))
}
