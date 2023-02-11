package hoard

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/koud-fi/pkg/blob"
	datastore "github.com/koud-fi/pkg/data"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/pk"
	"github.com/koud-fi/pkg/rx"
)

const Meta blob.Domain = "meta"

var _ blob.Storage = (*Hoard[any])(nil)

// Hoard is a content adressable storage.
type Hoard[T any] struct {
	file datastore.Table[File[T]]
	ref  datastore.Table[fileRef]
	data blob.Storage
	config

	// TODO: synchronization to make concurrent access less dodgy
}

type File[T any] struct {
	ID pk.UID `json:"id"`
	file.Attributes
	Metadata map[string]T `json:"metadata,omitempty"`
}

type fileRef struct {
	Ref string `json:"ref"`
	ID  pk.UID `json:"id"`
}

func New[T any](meta blob.SortedStorage, data blob.Storage, opt ...Option) *Hoard[T] {
	var c config
	for _, opt := range opt {
		opt(&c)
	}
	return &Hoard[T]{
		file: datastore.BlobsTable(meta, func(f File[T]) (blob.Ref, error) {
			return blob.NewRef(f.ID.Hex()), nil
		}),
		ref: datastore.BlobsTable(meta, func(r fileRef) (blob.Ref, error) {
			return blob.NewRef(r.Ref), nil
		}),
		data:   data,
		config: c,
	}
}

func (h *Hoard[T]) Get(ctx context.Context, ref blob.Ref) blob.Blob {
	return blob.Func(func() (io.ReadCloser, error) {
		switch ref.Domain() {
		case Meta:
			m, err := h.File(ctx, true, ref.Ref())
			if err != nil {
				return nil, err
			}
			if m.Ok() {
				return blob.Marshal(json.Marshal, m.Value()).Open()
			}
		default:
			if m, err := h.File(ctx, true, ref); err != nil {
				return nil, err
			} else if m.Ok() {
				return h.data.Get(ctx, blob.NewRef(m.Value().ID.Hex())).Open()
			}
		}
		return nil, os.ErrNotExist
	})
}

func (h *Hoard[T]) Set(ctx context.Context, ref blob.Ref, r io.Reader) error {
	switch ref.Domain() {
	case Meta:

		// ???

		panic("TODO")

	default:
		buf, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		var f File[T]

		if m, err := h.File(ctx, false, ref); err != nil {
			return err
		} else if !m.Ok() {
			if f.ID, err = pk.NewUID(buf, h.idSalt); err != nil {
				return err
			}
			if f.Attributes, err = file.ResolveAttrs(blob.FromBytes(buf), h.fileAttrOpts...); err != nil {
				return err
			}
			if err := h.data.Set(ctx, blob.NewRef(f.ID.Hex()), bytes.NewReader(buf)); err != nil {
				return err
			}
			for k, d := range f.Digest {
				ref := strings.ToLower(k) + ":" + d
				if _, err := h.ref.Put(ctx)(fileRef{Ref: ref, ID: f.ID}); err != nil {
					return err
				}
			}
		}
		if ref.Domain() != "" {
			if f.Metadata == nil {
				f.Metadata = make(map[string]T, 1)
			}
			key := ref.String()
			if _, ok := f.Metadata[key]; !ok {
				var init T
				f.Metadata[key] = init
			}
			if _, err := h.ref.Put(ctx)(fileRef{Ref: ref.String(), ID: f.ID}); err != nil {
				return err
			}
		}
		_, err = h.file.Put(ctx)(f)
		return err
	}
}

func (h *Hoard[T]) Delete(ctx context.Context, refs ...blob.Ref) error {

	// ???

	panic("TODO")
}

func (h Hoard[T]) File(ctx context.Context, resolve bool, ref blob.Ref) (rx.Maybe[File[T]], error) {
	switch ref.Domain() {
	case blob.Default:
		id, err := pk.ParseUID(ref.String())
		if err != nil {
			return rx.None[File[T]](), err
		}
		p, err := h.file.Get(ctx)(File[T]{ID: id})
		return p.Value(), err

	default:
		p, err := h.ref.Get(ctx)(fileRef{Ref: ref.String()})
		if err != nil {
			return rx.None[File[T]](), err
		}
		if m := p.Value(); m.Ok() {
			return h.File(ctx, resolve, blob.NewRef(m.Value().ID.Hex()))
		}
		if resolve && h.src != nil {
			buf, err := blob.Bytes(h.src.Get(ctx, ref))
			if err != nil {
				return rx.None[File[T]](), err
			}
			if err := h.Set(ctx, ref, bytes.NewReader(buf)); err != nil {
				return rx.None[File[T]](), err
			}
			return h.File(ctx, false, ref)
		}
		return rx.None[File[T]](), nil
	}
}
