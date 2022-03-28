package protoserver

import (
	"bytes"
	"context"
	"io/fs"
	"text/template"

	"github.com/koud-fi/pkg/pk"
)

func RegisterTemplate(s pk.Scheme, src fs.FS) {
	Register(s, templateFetcher{src: src})
}

func RegisterTemplateFS(fsys fs.FS) error {
	dir, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil
	}
	for _, d := range dir {

		// TODO: support single file templates

		if d.IsDir() {
			name := d.Name()
			src, err := fs.Sub(fsys, name)
			if err != nil {
				return err
			}
			RegisterTemplate(pk.Scheme(name), src)
		}
	}
	return nil
}

type templateFetcher struct {
	src fs.FS
}

func (f templateFetcher) Fetch(ctx context.Context, ref pk.Ref) (any, error) {

	// TODO: cache parsed templates

	t, err := template.ParseFS(f.src, "*")
	if err != nil {
		return nil, err
	}
	var (
		tc = templateContext{
			ctx: ctx,
			Key: ref.Key(),
		}
		buf  = bytes.NewBuffer(nil)
		name = ref.Params()
	)
	if name == "" {
		name = "index.tmpl"
	}
	if err := t.ExecuteTemplate(buf, name, tc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type templateContext struct {
	ctx context.Context
	Key string
}

func (tc templateContext) Fetch(refStr string) (any, error) {
	ref, err := pk.ParseRef(refStr)
	if err != nil {
		return nil, err
	}
	return Fetch(tc.ctx, ref)
}
