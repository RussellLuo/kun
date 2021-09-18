package oasv2

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var (
	tmplResponses = template.Must(template.New("responses").Parse(`
      produces:
        {{- range $contentType, $_ := .Responses.ContentTypes}}
        - {{$contentType}}
        {{- end}}
      responses:
        {{.Responses.Success.StatusCode}}:
          description: ""
          {{- if ne .Responses.Success.StatusCode 204}}
          schema:
			{{- if eq .Responses.Success.SchemaName "file"}}
            type: file
			{{- else}}
            $ref: "#/definitions/{{.Responses.Success.SchemaName}}"
			{{- end}}
          {{- end}}
        {{- range $statusCode, $response := .Responses.Failures}}
        {{$statusCode}}:
          description: ""
          schema:
            $ref: "#/definitions/{{$response.SchemaName}}"
        {{- end}}
`))

	funcs = map[string]interface{}{
		"basicJSONType": func(typ string) string {
			switch typ {
			case "bool":
				return "boolean"
			case "string":
				return "string"
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64":
				return "integer"
			case "float32", "float64":
				return "number"
			default:
				return ""
			}
		},
	}
	tmplDefinitions = template.Must(template.New("definitions").Funcs(funcs).Parse(`
definitions:
{{- range $name, $definition := .Definitions}}
  {{$name}}:
    {{- if eq $definition.Type "object"}}
    type: object

    {{- if $definition.ItemTypeOrProperties}}
    properties:
      {{- range $definition.ItemTypeOrProperties}}
      {{.Name}}:
        {{- if eq .Type.Kind "basic"}}
        type: {{.Type.Type}}

        {{- if .Type.Format}}
        format: {{.Type.Format}}
        {{- end -}} {{/* if .Type.Format */}}

        {{- else if eq .Type.Kind "object"}}
        $ref: "#/definitions/{{.Type.Type}}"

        {{- else if eq .Type.Kind "array"}}
        type: array
        items:
          {{- $basicJSONType := basicJSONType .Type.Type}}
          {{- if $basicJSONType}}
          type: {{$basicJSONType}}
          {{- else}}
          $ref: "#/definitions/{{.Type.Type}}"
          {{- end -}} {{/* if isBasic .Type.Type */}}
        {{- end -}} {{/* if eq .Type.Kind "basic" */}}

        {{- if .Type.Description}}
        description: {{.Type.Description}}
        {{- end}}

      {{- end -}} {{/* range $definition.ItemTypeOrProperties */}}
    {{- end -}} {{/* if $definition.ItemTypeOrProperties */}}

    {{- else if eq $definition.Type "array"}}
    type: array
    items:
      $ref: '#/definitions/{{$definition.ItemTypeOrProperties}}'
	{{- end -}} {{/* if eq $definition.Type "object" */}}
{{- end -}} {{/* range $name, $definition := .Definitions */}}
`))
)

func AddDefinition(defs map[string]Definition, name string, value reflect.Value) {
	parser := NewParser()
	parser.AddDefinition(name, value, false)
	for name, def := range parser.Definitions() {
		defs[name] = def
	}
}

func GetOASResponses(schema Schema, name string, statusCode int, body interface{}) OASResponses {
	resps := OASResponses{ContentTypes: map[string]bool{}, Failures: map[int]OASResponse{}}

	success := schema.SuccessResponse(name, statusCode, body)
	resps.Success = OASResponse{
		StatusCode: success.StatusCode,
		SchemaName: getSchemaName(success.ContentType, name+"Response"),
	}
	resps.ContentTypes[success.ContentType] = true

	failures := schema.FailureResponses(name)
	for _, failure := range failures {
		if _, ok := resps.Failures[failure.StatusCode]; ok {
			fmt.Printf("WARNING: Discard one response schema with %d for %s, since OAS-v2 does not support alternative schemas\n", failure.StatusCode, name)
		} else {
			resps.Failures[failure.StatusCode] = OASResponse{
				StatusCode: failure.StatusCode,
				SchemaName: name + "ResponseError" + strconv.Itoa(failure.StatusCode),
			}
		}
		resps.ContentTypes[failure.ContentType] = true
	}

	return resps
}

func AddResponseDefinitions(defs map[string]Definition, schema Schema, name string, statusCode int, body interface{}) {
	parser := NewParser()

	success := schema.SuccessResponse(name, statusCode, body)
	if success.Body != nil {
		parser.AddDefinition(name+"Response", reflect.ValueOf(success.Body), false)
	}

	failures := schema.FailureResponses(name)
	for _, failure := range failures {
		parser.AddDefinition(name+"ResponseError"+strconv.Itoa(failure.StatusCode), reflect.ValueOf(failure.Body), false)
	}

	for name, def := range parser.Definitions() {
		defs[name] = def
	}
}

func getSchemaName(contentType, defaultName string) string {
	if isMediaFile(contentType) {
		return "file"
	}
	return defaultName
}

func isMediaFile(contentType string) bool {
	if strings.HasPrefix(contentType, "image/png") {
		return true
	} else if strings.HasPrefix(contentType, "image/gif") {
		return true
	} else if strings.HasPrefix(contentType, "image/jpeg") {
		return true
	} else if strings.HasPrefix(contentType, "application/pdf") {
		return true
	}
	return false
}

func GenPaths(resps []OASResponses, paths string) string {
	var respStrings []interface{}
	for _, resp := range resps {
		data := struct {
			Responses OASResponses
		}{
			Responses: resp,
		}

		var buf bytes.Buffer
		if err := tmplResponses.Execute(&buf, data); err != nil {
			panic(err)
		}
		respStrings = append(respStrings, buf.String())
	}
	return fmt.Sprintf(paths, respStrings...)
}

func GenDefinitions(defs map[string]Definition) string {
	data := struct {
		Definitions map[string]Definition
	}{
		Definitions: defs,
	}

	var buf bytes.Buffer
	if err := tmplDefinitions.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}

type APIDocFunc func(schema Schema) string

func Handler(apiDocFn APIDocFunc, schema Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, apiDocFn(schema))
	}
}
