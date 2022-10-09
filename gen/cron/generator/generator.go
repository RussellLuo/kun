package generator

import (
	"github.com/RussellLuo/kun/gen/cron/parser"
	utilannotation "github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

var (
	template = utilannotation.FileHeader + `
{{- $srcPkgName := .Data.SrcPkgName}}
{{- $endpointPkgPrefix := .PkgInfo.EndpointPkgPrefix}}

package {{.PkgInfo.CurrentPkgName}}

import (
	{{- range .Data.Imports}}
	{{.ImportString}}
	{{- end}}
	"github.com/RussellLuo/micron"

	{{- if .PkgInfo.EndpointPkgPath}}
	"{{.PkgInfo.EndpointPkgPath}}"
	{{- end}}
)

func NewCronJobs(svc {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}}) []micron.Job {
	return []micron.Job{
		{{range .Spec.Operations -}}
		{
			Name: "{{getJobName .GoMethodName}}",
			Expr: "{{getJobExpr .GoMethodName}}",
			Handler: svc.{{.Name}},
		},
		{{end -}} {{/* range .Spec.Operations */ -}}
	}
}
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

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, ifaceData *ifacetool.Data, cronSpec map[string]*parser.Job, spec *openapi.Specification) (*generator.File, error) {
	data := struct {
		PkgInfo *generator.PkgInfo
		Data    *ifacetool.Data
		Spec    *openapi.Specification
	}{
		PkgInfo: pkgInfo,
		Data:    ifaceData,
		Spec:    spec,
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"getJobName": func(methodName string) string {
				return cronSpec[methodName].Name
			},
			"getJobExpr": func(methodName string) string {
				return cronSpec[methodName].Expr
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "cron.go",
	})
}
