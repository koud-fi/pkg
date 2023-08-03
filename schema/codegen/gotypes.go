package codegen

import (
	"html/template"
	"io"

	"github.com/koud-fi/pkg/schema"
)

var goTypeTemplate = template.Must(template.New("go-type").
	Funcs(map[string]any{
		"typeOf": func(t schema.Type) string {
			switch t.Type {
			case schema.String:
				return "string"
			case schema.Number:
				return "float64"
			case schema.Integer:
				return "int"
			case schema.Object:
				return "map[string]any" // TODO: use a struct
			case schema.Array:
				return "[]any" // TODO: item type
			case schema.Boolean:
				return "bool"
			default:
				return "any"
			}
		},
		"tags": func(t schema.Type, key string) string {
			return "TODO"
		},
	}).
	Parse(`
package codegen

// This code is generated by go generate.
// DO NOT EDIT BY HAND!

{{range $typeName, $typeDef := .Definitions -}}
type {{$typeName}} struct {
{{- range $fieldName, $fieldDef := .Properties}}
	{{$fieldName}} {{typeOf $fieldDef -}}
{{- end }}
}
{{- end }}
`))

func GoTypes(w io.Writer, s schema.Schema) error {
	return goTypeTemplate.Execute(w, s)
}
