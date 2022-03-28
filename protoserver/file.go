package protoserver

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/koud-fi/pkg/pk"
	"github.com/koud-fi/pkg/schema"
)

const FileScheme = "file"

func RegisterFile() {
	Register(FileScheme,
		TransformFetcher(func(ctx context.Context, key string, ref pk.Ref) (any, error) {

			// TODO: caching of non-local files

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
			return f.Open()
		}))
}
