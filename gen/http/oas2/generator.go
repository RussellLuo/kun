package oas2

import (
	"fmt"
	"strings"

	"github.com/RussellLuo/kun/gen/http/parser/annotation"
	utilannotation "github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/caseconv"
)

var (
	template = utilannotation.FileHeader + `
package {{.PkgInfo.CurrentPkgName}}

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"text/template"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/RussellLuo/kun/pkg/oas2"

	{{- if .PkgInfo.EndpointPkgPath}}
	"{{.PkgInfo.EndpointPkgPath}}"
	{{- end}}
)

var (
	base = ` + "`" + `swagger: "2.0"
info:
  title: "{{.Spec.Metadata.Title}}"
  version: "{{.Spec.Metadata.Version}}"
  description: "{{.Spec.Metadata.Description}}"
  license:
    name: "MIT"
host: "example.com"
basePath: "{{.Spec.Metadata.BasePath}}"
schemes:
  - "https"
consumes:
  - "application/json"
produces:
  - "application/json"
` + "`" + `

{{- $defaultTags := .Spec.Metadata.DefaultTags}}
{{- $operationsGroupByPattern := operationsGroupByPattern .Spec.Operations}}

	paths = ` + "`" + `
paths:
{{- range $operationsGroupByPattern}}
  {{.Pattern}}:

  {{- range .Operations}}
  {{- $nonCtxParams := nonCtxParams .Request.Params}}
    {{lower .Method}}:
      description: "{{.Description}}"
      summary: "{{.Description}}"
      operationId: "{{.Name}}"
      {{- $tags := getTags .Tags $defaultTags}}
      {{- if $tags}}
      tags:
        {{- range $tags}}
        - {{.}}
        {{- end -}} {{/* range $tags */}}
      {{- end -}} {{/* if $tags */}}

      {{- $nonCtxNonBodyParams := nonBodyParams $nonCtxParams}}
      {{- if $nonCtxParams}}
      parameters:
        {{- range $nonCtxNonBodyParams}}
        - name: {{.Alias}}
          in: {{.In}}
          required: {{.Required}}
          {{- $type := typeName .AliasType}}
          type: {{$type.Type}}
          {{- if $type.ItemType}}
          items:
            type: {{$type.ItemType}}
          {{- end}}
          description: "{{.Description}}"
        {{- end -}} {{/* range $nonCtxNonBodyParams */}}

        {{- $bodyParams := bodyParams $nonCtxParams}}
        {{- if $bodyParams}}
        - name: body
          in: body
          schema:
            $ref: "#/definitions/{{.GoMethodName}}RequestBody"
        {{- end}}
      {{- end -}} {{/* if $nonCtxParams */}}
      %s
  {{- end -}} {{/* range .Operations */}}

{{- end -}} {{/* range $operationsGroupByPattern */}}
` + "`" + `
)

func getResponses(schema oas2.Schema) []oas2.OASResponses {
	return []oas2.OASResponses{
		{{- range $operationsGroupByPattern}}
		{{- range .Operations}}
		oas2.GetOASResponses(schema, "{{.GoMethodName}}", {{.SuccessResponse.StatusCode}}, {{endpointPrefix .GoMethodName}}Response{}),
		{{- end}} {{/* range .Operations */}}
		{{- end}} {{/* range $operationsGroupByPattern */}}
	}
}

func getDefinitions(schema oas2.Schema) map[string]oas2.Definition {
	defs := make(map[string]oas2.Definition)

	{{range .Spec.Operations -}}

	{{- $nonCtxParams := nonCtxParams .Request.Params}}
	{{- $bodyParams := bodyParams $nonCtxParams}}
	{{- $bodyField := getBodyField .Request.BodyField}}
	{{- if $bodyField}}
	oas2.AddDefinition(defs, "{{.GoMethodName}}RequestBody", reflect.ValueOf(({{endpointPrefix .GoMethodName}}Request{}).{{title $bodyField}}))
	{{- else if $bodyParams}}
	oas2.AddDefinition(defs, "{{.GoMethodName}}RequestBody", reflect.ValueOf(&struct{
		{{- range $bodyParams}}
		{{title .Name}} {{.Type}} {{addTag .Alias .Type .Description .Required }}
		{{- end}} {{/* range $bodyParams */}}
	}{}))
	{{- end}} {{/* if $bodyField */}}
	oas2.AddResponseDefinitions(defs, schema, "{{.GoMethodName}}", {{.SuccessResponse.StatusCode}}, ({{endpointPrefix .GoMethodName}}Response{}).Body())

    {{end -}} {{/* range .Spec.Operations */}}

	return defs
}

func OASv2APIDoc(schema oas2.Schema) string {
	resps := getResponses(schema)
	paths := oas2.GenPaths(resps, paths)

	defs := getDefinitions(schema)
	definitions := oas2.GenDefinitions(defs)

	return base + paths + definitions
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

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, spec *openapi.Specification) (*generator.File, error) {
	data := struct {
		PkgInfo *generator.PkgInfo
		Spec    *openapi.Specification
	}{
		PkgInfo: pkgInfo,
		Spec:    spec,
	}

	type OperationsPerPattern struct {
		Pattern    string
		Operations []*openapi.Operation
	}

	type ParamType struct {
		Type     string
		ItemType string
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"title": caseconv.UpperFirst,
			"lower": strings.ToLower,
			"operationsGroupByPattern": func(ops []*openapi.Operation) (outOps []*OperationsPerPattern) {
				var opp *OperationsPerPattern
				var ok bool

				patternToOps := make(map[string]*OperationsPerPattern)

				for _, op := range ops {
					opp, ok = patternToOps[op.Pattern]
					if !ok {
						opp = &OperationsPerPattern{Pattern: op.Pattern}
						outOps = append(outOps, opp)
						patternToOps[op.Pattern] = opp
					}
					op.Description = strings.TrimSpace(strings.TrimPrefix(op.Description, op.Name))
					opp.Operations = append(opp.Operations, op)
				}

				return
			},
			"typeName": func(typ string) ParamType {
				basicOASType := func(n string) string {
					switch n {
					case "int", "int8", "int16", "int32", "int64",
						"uint", "uint8", "uint16", "uint32", "uint64":
						return "integer"
					case "float32", "float64":
						return "number"
					case "string":
						return "string"
					case "bool":
						return "boolean"
					default:
						return "string"
					}
				}

				isBasicGoType := func(t string) bool {
					switch t {
					case "bool", "string",
						"int", "int8", "int16", "int32", "int64",
						"uint", "uint8", "uint16", "uint32", "uint64",
						"float32", "float64":
						return true
					default:
						return false
					}
				}

				// Cases are more complicated regarding Argument aggregation,
				// where typ (as the same as argument type) usually does not
				// match the actual type of the corresponding request parameter.
				//
				// | Method argument | Request parameter    |
				// | --------------- | -------------------- |
				// | basic or slice  | string (UNSUPPORTED) |
				// | struct or other | string               |
				//
				// Adding one more optional tag `type` in the @kun(param)
				// annotation may be a better solution?

				switch {
				case isBasicGoType(typ):
					return ParamType{Type: basicOASType(typ)}
				case strings.HasPrefix(typ, "[]") && isBasicGoType(typ[2:]):
					itemType := typ[2:]
					return ParamType{Type: "array", ItemType: basicOASType(itemType)}
				default:
					return ParamType{Type: "string"}
				}
			},
			"getTags": func(opTags, defaultTags []string) []string {
				if len(opTags) > 0 {
					return opTags
				}
				return defaultTags
			},
			"nonCtxParams": func(params []*openapi.Param) (out []*openapi.Param) {
				for _, p := range params {
					if p.Type != "context.Context" && p.In != openapi.InRequest {
						// Parameters in `request` have no relationship with OAS.
						out = append(out, p)
					}
				}
				return
			},
			"nonBodyParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In != openapi.InBody {
						out = append(out, p)
					}
				}
				return
			},
			"bodyParams": func(in []*openapi.Param) (out []*openapi.Param) {
				for _, p := range in {
					if p.In == openapi.InBody {
						out = append(out, p)
					}
				}
				return
			},
			"endpointPrefix": func(name string) string {
				fullName := pkgInfo.EndpointPkgPrefix + name
				if g.opts.SchemaPtr {
					return "&" + fullName
				}
				return fullName
			},
			"addTag": func(name, typ, description string, required bool) string {
				if g.opts.SchemaTag == "" {
					return ""
				}

				if typ == "error" {
					name = "-"
				}

				tag := fmt.Sprintf(`%s:"%s"`, g.opts.SchemaTag, name)

				kunTag := func(descr string, required bool) string {
					var content []string
					if descr != `` {
						content = append(content, fmt.Sprintf(`descr=%s`, descr))
					}
					if required {
						content = append(content, fmt.Sprintf(`required=%v`, required))
					}
					if len(content) == 0 {
						return ``
					}
					return fmt.Sprintf(`kun:"%s"`, strings.Join(content, ` `))
				}(description, required)

				if kunTag != `` {
					tag = tag + ` ` + kunTag
				}
				return "`" + tag + "`"
			},
			"getBodyField": func(name string) string {
				if name != "" && name != annotation.OptionNoBody {
					return name
				}
				return ""
			},
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "oas2.go",
	})
}
