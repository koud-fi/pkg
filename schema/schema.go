package schema

import (
	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/pk"
)

const (
	MasterFileKey = "master"

	PageNext     LinkType = "page:next"
	PagePrevious LinkType = "page:prev"
)

type RefNode interface {
	Ref() (pk.Ref, error)
}

type ParentNode interface {
	Children() ([]RefNode, error)
}

type LinkNode interface {
	Links(current pk.Ref) (map[LinkType]pk.Ref, error)
}

type LinkType string

type TagNode interface {
	Tags() (map[string]string, error)
}

type FileNode interface {
	Files() (map[string]File, error)
}

type File interface {
	blob.Blob
	// TODO: attributes interface
}
