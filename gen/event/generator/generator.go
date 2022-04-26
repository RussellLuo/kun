package generator

import (
	"fmt"
	"strings"

	"github.com/RussellLuo/kun/gen/event/parser"
	utilannotation "github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/caseconv"
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

	{{- if .PkgInfo.EndpointPkgPath}}
	"{{.PkgInfo.EndpointPkgPath}}"
	{{- end}}
)

func NewEventHandler(svc {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}}, codecs eventcodec.Codecs) eventpubsub.Handler {
	var codec eventcodec.Codec
	handlerSet := eventpubsub.NewHandlerSet()

	{{range .Spec.Operations -}}
	codec = codecs.EncodeDecoder("{{.GoMethodName}}")
	handlerSet.Add("{{getEventType .GoMethodName}}", eventpubsub.NewSubscriber(
		{{$endpointPkgPrefix}}MakeEndpointOf{{.GoMethodName}}(svc),
		decode{{.Name}}Input(codec),
	))

	{{end -}} {{/* range .Spec.Operations */ -}}
	return handlerSet
}

{{- range .Spec.Operations}}

{{- $nonCtxParams := nonCtxParams .Request.Params}}
{{- $methodHasNonCtxParams := methodHasNonCtxParams .GoMethodName}}
{{- $hasBodyParams := hasBodyParams $nonCtxParams}}
{{- $dataField := getDataField}}

func decode{{.Name}}Input(codec eventcodec.Codec) eventpubsub.DecodeInputFunc {
	return func(_ context.Context, event eventpubsub.Event) (interface{}, error) {
		{{- if $methodHasNonCtxParams}}
		var input {{$endpointPkgPrefix}}{{.GoMethodName}}Request

		{{end -}}

		{{if $dataField -}}
		if err := codec.Decode(event.Data(), &input.{{title $dataField}}); err != nil {
			return nil, err
		}
		{{else if $hasBodyParams -}}
		if err := codec.Decode(event.Data(), &input); err != nil {
			return nil, err
		}
		{{end -}}

		{{- if $methodHasNonCtxParams}}

		return {{addAmpersand "input"}}, nil
		{{- else -}}
		return nil, nil
		{{- end}} {{/* End of if $methodHasNonCtxParams */}}
	}
}

{{- end}}
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

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, ifaceData *ifacetool.Data, eventInfo *parser.EventInfo, spec *openapi.Specification) (*generator.File, error) {
	data := struct {
		PkgInfo *generator.PkgInfo
		Data    *ifacetool.Data
		Spec    *openapi.Specification
		Opts    *Options
	}{
		PkgInfo: pkgInfo,
		Data:    ifaceData,
		Spec:    spec,
		Opts:    g.opts,
	}

	methodMap := make(map[string]*ifacetool.Method)
	for _, method := range ifaceData.Methods {
		methodMap[method.Name] = method
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"title":      strings.Title,
			"lowerFirst": caseconv.LowerFirst,
			"addAmpersand": func(name string) string {
				if g.opts.SchemaPtr {
					return "&" + name
				}
				return name
			},
			"hasBodyParams": func(in []*openapi.Param) bool {
				for _, p := range in {
					if p.In == openapi.InBody {
						return true
					}
				}
				return false
			},
			"nonCtxParams": func(params []*openapi.Param) (out []*openapi.Param) {
				for _, p := range params {
					if p.Type != "context.Context" {
						out = append(out, p)
					}
				}
				return
			},
			"methodHasNonCtxParams": func(methodName string) bool {
				method, ok := methodMap[methodName]
				if !ok {
					panic(fmt.Errorf("no method named %q", methodName))
				}

				for _, p := range method.Params {
					if p.TypeString != "context.Context" {
						return true
					}
				}
				return false
			},
			"getDataField": func() string {
				return eventInfo.DataField
			},
			"getEventType": func(methodName string) string {
				return eventInfo.Types[methodName]
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "event.go",
	})
}
