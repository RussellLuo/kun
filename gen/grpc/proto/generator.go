package proto

import (
	"fmt"

	"github.com/RussellLuo/kun/gen/grpc/parser"
	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/pkgtool"
)

var (
	template = annotation.FileHeader + `
syntax = "proto3";

option go_package = "{{.PkgPath}}";

package {{.PkgName}};

{{range .Service.Descriptions -}}
{{.}}
{{- end}}
service {{.Service.Name}} {
  {{- range .Service.RPCs}}
  {{- range .Descriptions}}
  {{.}}
  {{- end}} {{/* range .Description */}}
  rpc {{.Name}} ({{.Request.Name}}) returns ({{.Response.Name}}) {}
  {{- end}} {{/* range .Service.RPCs */}}
}

{{- range .Service.RPCs}}

// The request message of {{.Name}}.
message {{.Request.Name}} {
  {{- range .Request.Fields}}
  {{fullTypeName .Type}} {{snakeCase .Name}} = {{.Num}};
  {{- end}} {{/* range .Request.Fields */}}
}

// The response message of {{.Name}}.
message {{.Response.Name}} {
  {{- range .Response.Fields}}
  {{fullTypeName .Type}} {{snakeCase .Name}} = {{.Num}};
  {{- end}} {{/* range .Response.Fields */}}
}
{{- end}} {{/* range .Service.RPCs */}}

{{- range .Messages}}

{{if .Fields -}}
message {{.Name}} {
  {{- range .Fields}}
  {{fullTypeName .Type}} {{snakeCase .Name}} = {{.Num}};
  {{- end}} {{/* range .Fields */}}
}
{{- end}}{{/* if .Fields */}}

{{- end}} {{/* range .Messages */}}
`
)

type Options struct {
	SchemaPtr bool
	SchemaTag string
	Formatted bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(outDir string, ifaceData *ifacetool.Data, service *parser.Service) (*generator.File, error) {
	data := struct {
		PkgPath  string
		PkgName  string
		Service  *parser.Service
		Messages map[string]*parser.Type
	}{
		PkgPath:  pkgtool.PkgPathFromDir(outDir),
		PkgName:  pkgtool.PkgNameFromDir(outDir),
		Service:  service,
		Messages: getMessages(service),
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"snakeCase": caseconv.ToSnakeCase,
			"fullTypeName": func(typ *parser.Type) string {
				name := typ.Name

				if typ.MapKey != "" {
					name = fmt.Sprintf("map<%s, %s>", typ.MapKey, typ.Name)
				}

				if typ.Repeated {
					name = "repeated " + name
				}

				return name
			},
		},
		TargetFileName: ifaceData.SrcPkgName + ".proto",
	})
}

func getMessages(s *parser.Service) map[string]*parser.Type {
	// TODO: Import another .proto's definitions, if necessary, for best reusability.

	types := make(map[string]*parser.Type)

	addTypes := func(fields []*parser.Field) {
		for _, f := range fields {
			for _, t := range f.Type.Squash() {
				types[t.Name] = t
			}
		}
	}

	for _, rpc := range s.RPCs {
		addTypes(rpc.Request.Fields)
		addTypes(rpc.Response.Fields)
	}

	return types
}
