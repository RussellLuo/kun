package chi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RussellLuo/kok/gen/util/annotation"
	"github.com/RussellLuo/kok/gen/util/generator"
	"github.com/RussellLuo/kok/gen/util/openapi"
	"github.com/RussellLuo/kok/pkg/caseconv"
	"github.com/RussellLuo/kok/pkg/ifacetool"
)

var (
	template = annotation.FileHeader + `
{{- $srcPkgName := .Data.SrcPkgName}}
{{- $endpointPkgPrefix := .PkgInfo.EndpointPkgPrefix}}
{{- $enableTracing := .Opts.EnableTracing}}

package {{.PkgInfo.CurrentPkgName}}

import (
	"encoding/json"
	"net/http"
	"strconv"

	{{- if $enableTracing}}
	"github.com/RussellLuo/kok/pkg/trace/xnet"
	{{- end}}
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/RussellLuo/kok/pkg/httpoption2"

	{{- range .Data.Imports}}
	{{.ImportString}}
	{{- end}}

	{{- if .PkgInfo.EndpointPkgPath}}
	"{{.PkgInfo.EndpointPkgPath}}"
	{{- end}}
)

func NewHTTPRouter(svc {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}}, codecs httpcodec.Codecs, opts ...httpoption.Option) chi.Router {
	r := chi.NewRouter()
	options := httpoption.NewOptions(opts...)

	{{if $enableTracing -}}
	contextor := xnet.NewContextor()
	r.Method("PUT", "/trace", xnet.HTTPHandler(contextor))
	{{- end}}

	r.Method("GET", "{{.Spec.Metadata.DocsPath}}", oas2.Handler(OASv2APIDoc, options.ResponseSchema()))

	var codec httpcodec.Codec
	var validator httpoption.Validator
	var kitOptions []kithttp.ServerOption

	{{- range .Spec.Operations}}

	codec = codecs.EncodeDecoder("{{.Name}}")
	validator = options.RequestValidator("{{.Name}}")
	r.Method(
		"{{.Method}}", "{{.Pattern}}",
		kithttp.NewServer(
			{{$endpointPkgPrefix}}MakeEndpointOf{{.Name}}(svc),
			decode{{.Name}}Request(codec, validator),
			httpcodec.MakeResponseEncoder(codec, {{getStatusCode .SuccessResponse.StatusCode .Name}}),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
				{{- if $enableTracing}}
				kithttp.ServerBefore(contextor.HTTPToContext("{{$srcPkgName}}", "{{.Name}}")),
				{{- end}}
			)...,
		),
	)
	{{- end}}

	return r
}

func NewHTTPRouterWithOAS(svc {{$.Data.SrcPkgQualifier}}{{$.Data.InterfaceName}}, codecs httpcodec.Codecs, schema oas2.Schema) chi.Router {
	return NewHTTPRouter(svc, codecs, httpoption.ResponseSchema(schema))
}

{{- range .Spec.Operations}}

{{- $nonCtxParams := nonCtxParams .Request.Params}}
{{- $nonBodyParamsGroupByName := nonBodyParamsGroupByName $nonCtxParams}}
{{- $hasBodyParams := hasBodyParams $nonCtxParams}}
{{- $bodyField := getBodyField .Request.BodyField}}

func decode{{.Name}}Request(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		{{- if $nonCtxParams}}
		var _req {{$endpointPkgPrefix}}{{.Name}}Request

		{{end -}}

		{{if $bodyField -}}
		if err := codec.DecodeRequestBody(r, &_req.{{title $bodyField}}); err != nil {
			return nil, err
		}
		{{else if $hasBodyParams -}}
		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}
		{{end -}}

		{{- range $nonBodyParamsGroupByName}}

		{{- if .Aggregation}}
		{{.Name}} := map[string][]string{
			{{- range .Properties}}
			"{{.In}}.{{.Alias}}": {{extractParam .}},
			{{- end}}
		}
		if err := codec.DecodeRequestParams("{{.Name}}", {{.Name}}, {{paramVar "_req" .}}); err != nil {
			return nil, err
		}

		{{- else}} {{/* if .Aggregation */}}
		{{- $property := index .Properties 0}}
		{{.Name}} := {{$property | extractParam}}
		{{- if $property.Required}}
		if err := codec.DecodeRequestParam("{{.Name}}", {{.Name}}, {{paramVar "_req" .}}); err != nil {
			return nil, err
		}
		{{- else}} {{/* if $property.Required */}}
		if len({{.Name}}) > 0 {
			if err := codec.DecodeRequestParam("{{.Name}}", {{.Name}}, {{paramVar "_req" .}}); err != nil {
				return nil, err
			}
		}
		{{- end}} {{/* if $property.Required */}}

		{{- end}} {{/* if .Aggregation */}}

		{{end -}} {{/* range $nonBodyParamsGroupByName */}}

		{{- if $nonCtxParams}}

		if err := validator.Validate({{addAmpersand "_req"}}); err != nil {
			return nil, err
		}

		return {{addAmpersand "_req"}}, nil
		{{- else -}}
		return nil, nil
		{{- end}} {{/* End of if $nonCtxParams */}}
	}
}

{{- end}}
`
)

type RequestField struct {
	Name  string
	Value string
}

type Server struct {
	Service     interface{}
	NewEndpoint interface{}
	Request     interface{}
	Response    interface{}
}

type Options struct {
	SchemaPtr     bool
	SchemaTag     string
	Formatted     bool
	EnableTracing bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, ifaceData *ifacetool.Data, spec *openapi.Specification) (*generator.File, error) {
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

	type ParamProperty struct {
		In       string
		Alias    string
		Required bool
	}

	type ParamsGroupByName struct {
		Name        string
		IsBlank     bool
		Aggregation bool
		Properties  []ParamProperty
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
			"extractParam": func(param *ParamProperty) string {
				switch param.In {
				case openapi.InPath:
					return fmt.Sprintf(`[]string{chi.URLParam(r, "%s")}`, param.Alias)
				case openapi.InQuery:
					return fmt.Sprintf(`r.URL.Query()["%s"]`, param.Alias)
				case openapi.InHeader:
					return fmt.Sprintf(`r.Header.Values("%s")`, param.Alias)
				case openapi.InRequest:
					return fmt.Sprintf(`[]string{r.%s}`, param.Alias)
				default:
					panic(fmt.Errorf("param.In `%s` not supported", param.In))
				}
			},
			"nonBodyParamsGroupByName": func(in []*openapi.Param) (out []*ParamsGroupByName) {
				var names []string
				params := make(map[string]*ParamsGroupByName)

				for _, p := range in {
					if p.In != openapi.InBody {
						grouped, ok := params[p.Name]
						if !ok {
							grouped = &ParamsGroupByName{Name: p.Name, IsBlank: p.IsBlank}

							names = append(names, p.Name)
							params[p.Name] = grouped
						}
						grouped.Properties = append(grouped.Properties, ParamProperty{
							In:       p.In,
							Alias:    p.Alias,
							Required: p.Required,
						})
					}
				}

				for _, name := range names {
					p := params[name]
					if len(p.Properties) > 1 {
						p.Aggregation = true
					}
					out = append(out, p)
				}
				return
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
			"getStatusCode": func(givenStatusCode int, name string) int {
				method, ok := methodMap[name]
				if !ok {
					panic(fmt.Errorf("no method named %q", name))
				}

				if len(method.Returns) > 0 {
					// Use the given status code, since the corresponding
					// method is a fruitful function.
					return givenStatusCode
				}

				if givenStatusCode == http.StatusOK {
					fmt.Printf("NOTE: statusCode is changed to be 204, since method %q returns no result\n", name)
					return http.StatusNoContent
				}

				if givenStatusCode != http.StatusNoContent {
					panic(fmt.Errorf("statusCode must be 204, since method %q returns no result", name))
				}
				return givenStatusCode
			},
			"getBodyField": func(name string) string {
				if name != "" && name != openapi.OptionNoBody {
					return name
				}
				return ""
			},
			"paramVar": func(reqVar string, param *ParamsGroupByName) string {
				if param.IsBlank {
					return "nil"
				}
				return fmt.Sprintf("&%s.%s", reqVar, strings.Title(param.Name))
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "http.go",
	})
}
