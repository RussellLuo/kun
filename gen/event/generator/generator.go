package generator

import (
	"fmt"

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
	"github.com/RussellLuo/kun/pkg/eventcodec"
	"github.com/RussellLuo/kun/pkg/eventpubsub"

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
		{{end -}} {{/* if $dataField */}}

		{{- if $methodHasNonCtxParams}}

		return {{addAmpersand "input"}}, nil
		{{- else -}}
		return nil, nil
		{{- end}} {{/* if $methodHasNonCtxParams */}}
	}
}

{{- end}} {{/* range .Spec.Operations */}}

// EventPublisher implements {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}} on the publisher side.
//
// EventPublisher should only be used in limited scenarios where only one subscriber
// is involved and the publisher depends on the interface provided by the subscriber.
//
// In typical use cases of the publish-subscribe pattern - many subscribers are
// involved and the publisher knows nothing about the subscribers - you should
// just send the event in the way it should be.
type EventPublisher struct {
	publisher eventpubsub.Publisher
	codecs    eventcodec.Codecs
}

func NewEventPublisher(publisher eventpubsub.Publisher, codecs eventcodec.Codecs) *EventPublisher {
	return &EventPublisher{
		publisher: publisher,
		codecs:    codecs,
	}
}

{{- range .Spec.Operations}}

{{- $nonCtxParams := nonCtxParams .Request.Params}}
{{- $methodHasNonCtxParams := methodHasNonCtxParams .GoMethodName}}
{{- $dataField := getDataField}}
{{- $method := getMethod .GoMethodName}}

func (p *EventPublisher) {{$method.Name}}({{$method.ArgList}}) {{$method.ReturnArgNamedValueList}} {
	{{- if $nonCtxParams}}
	codec := p.codecs.EncodeDecoder("{{.GoMethodName}}")

	{{if $dataField -}}
	_data, err := codec.Encode({{$dataField}})
	{{- else}}
	_data, err := codec.Encode({{addAmpersand ""}}{{$endpointPkgPrefix}}{{.GoMethodName}}Request{
		{{- range $nonCtxParams}}
		{{title .Name}}: {{.Name}},
		{{- end}}
	})
	{{- end}} {{/* if $dataField */}}
	if err != nil {
		return err
	}
	{{- end}} {{/* if $nonCtxParams */}}

	{{- if $nonCtxParams}}

	return p.publisher.Publish({{getCtxArg .GoMethodName}}, "{{getEventType .GoMethodName}}", _data)
	{{- else}}
	return p.publisher.Publish({{getCtxArg .GoMethodName}}, "{{getEventType .GoMethodName}}", nil)
	{{- end}} {{/* if $nonCtxParams */}}
}

{{- end}} {{/* range .Spec.Operations */}}
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
			"title":      caseconv.UpperFirst,
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
			"getMethod": func(methodName string) *ifacetool.Method {
				method, ok := methodMap[methodName]
				if !ok {
					panic(fmt.Errorf("no method named %q", methodName))
				}
				return method
			},
			"getCtxArg": func(methodName string) string {
				method, ok := methodMap[methodName]
				if !ok {
					panic(fmt.Errorf("no method named %q", methodName))
				}

				for _, p := range method.Params {
					if p.TypeString == "context.Context" {
						return p.Name
					}
				}
				return "context.Background()"
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "event.go",
	})
}
