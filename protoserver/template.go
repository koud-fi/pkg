package protoserver

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"

	"github.com/koud-fi/pkg/pk"
)

func RegisterTemplate(s pk.Scheme, src fs.FS) {
	Register(s, templateFetcher{src: src})
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
	tc := templateContext{ctx: ctx}
	if ref.Key() != "" {
		tc.Ref, err = pk.ParseRef(ref.Key())
		if err != nil {
			return nil, err
		}
	}
	var (
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
	Ref pk.Ref
}

func (tc templateContext) Fetch(ref pk.Ref) (any, error) {
	return Fetch(tc.ctx, ref)
}
