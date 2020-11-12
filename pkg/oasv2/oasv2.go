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

	tmplDefinitions = template.Must(template.New("definitions").Parse(`
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
        {{- end}} {{/* if .Type.Format */}}

        {{- else if eq .Type.Kind "object"}}
        $ref: "#/definitions/{{.Type.Type}}"

        {{- else if eq .Type.Kind "array"}}
        type: array
        items:
          $ref: "#/definitions/{{.Type.Type}}"
        {{- end}} {{/*if eq .Type.Kind "basic" */}}

      {{- end}} {{/* range $definition.ItemTypeOrProperties */}}
    {{- end}} {{/* if $definition.ItemTypeOrProperties */}}

    {{- else if eq $definition.Type "array"}}
    type: array
    items:
      $ref: '#/definitions/{{$definition.ItemTypeOrProperties}}'
	{{- end}} {{/* if eq $definition.Type "object" */}}
{{- end}} {{/* range $name, $definition := .Definitions */}}
`))
)

type JSONType struct {
	Kind   string
	Type   string
	Format string
}

type ItemType JSONType

type Property struct {
	Name string
	Type JSONType
}

type Definition struct {
	Type                 string
	ItemTypeOrProperties interface{}
}

type OASResponse struct {
	StatusCode int
	SchemaName string
}

type OASResponses struct {
	ContentTypes map[string]bool
	Success      OASResponse
	Failures     map[int]OASResponse
}

func AddDefinition(defs map[string]Definition, name string, value reflect.Value) {
	if _, ok := defs[name]; ok {
		// Ignore duplicated definitions implicitly.
		return
	}

	switch value.Kind() {
	case reflect.Struct:
		var properties []Property

		structType := value.Type()
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldName := field.Name
			jsonTag := field.Tag.Get("json")
			jsonName := strings.SplitN(jsonTag, ",", 2)[0]
			if jsonName != "" {
				if jsonName == "-" {
					continue
				}
				fieldName = jsonName
			}

			fieldValue := addSubDefinition(defs, fieldName, value.Field(i))

			properties = append(properties, Property{
				Name: fieldName,
				Type: getJSONType(fieldValue.Type(), fieldName),
			})
		}

		defs[name] = Definition{
			Type:                 "object",
			ItemTypeOrProperties: properties,
		}
	case reflect.Map:
		var properties []Property

		valueType := value.Type()
		if kind := valueType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
			panic(fmt.Errorf(
				"'%s' needs a map with string keys, has '%s' keys",
				name, valueType.Key().Kind()))
		}

		for _, key := range value.MapKeys() {
			keyString := key.String()
			keyValue := addSubDefinition(defs, keyString, value.MapIndex(key))

			properties = append(properties, Property{
				Name: keyString,
				Type: getJSONType(keyValue.Type(), keyString),
			})
		}

		defs[name] = Definition{
			Type:                 "object",
			ItemTypeOrProperties: properties,
		}
	case reflect.Slice, reflect.Array:
		elemType := value.Type().Elem()

		if elemType.Kind() != reflect.Struct {
			panic(fmt.Errorf("only struct slice or array is supported, but got %v", elemType.String()))
		}

		structValue := reflect.New(elemType).Elem()
		AddDefinition(defs, elemType.Name(), structValue)

	case reflect.Ptr:
		elemType := value.Type().Elem()
		if elemType.Kind() != reflect.Struct {
			panic(fmt.Errorf("only struct pointer is supported, but got %v", elemType.String()))
		}

		elem := value.Elem()
		if !elem.IsValid() {
			elem = reflect.New(elemType).Elem()
		}
		AddDefinition(defs, elemType.Name(), elem)

	default:
		panic(fmt.Errorf("unsupported type %s", value.Kind()))
	}
}

func addSubDefinition(defs map[string]Definition, name string, value reflect.Value) reflect.Value {
	typeName := value.Type().Name()
	if typeName == "" {
		typeName = strings.Title(name)
	}

	switch value.Kind() {
	case reflect.Map:
		AddDefinition(defs, typeName, value)
	case reflect.Slice, reflect.Array:
		switch value.Type().Elem().Kind() {
		case reflect.Bool, reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		default:
			AddDefinition(defs, typeName, value)
		}
	case reflect.Ptr:
		switch value.Elem().Kind() {
		case reflect.Bool, reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		default:
			AddDefinition(defs, typeName, value)
		}
	case reflect.Struct:
		AddDefinition(defs, typeName, value)
	case reflect.Interface:
		value = addSubDefinition(defs, typeName, value.Elem())
	}

	return value
}

func getJSONType(typ reflect.Type, name string) JSONType {
	switch typ.Kind() {
	case reflect.Bool:
		return JSONType{Kind: "basic", Type: "boolean"}
	case reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return JSONType{Kind: "basic", Type: "integer", Format: "int32"}
	case reflect.Int, reflect.Int64,
		reflect.Uint, reflect.Uint64, reflect.Uintptr:
		return JSONType{Kind: "basic", Type: "integer", Format: "int64"}
	case reflect.Float32:
		return JSONType{Kind: "basic", Type: "number", Format: "float"}
	case reflect.Float64:
		return JSONType{Kind: "basic", Type: "number", Format: "double"}
	case reflect.String:
		return JSONType{Kind: "basic", Type: "string"}
	case reflect.Struct:
		return JSONType{Kind: "object", Type: typ.Name()}
	case reflect.Map:
		return JSONType{Kind: "object", Type: strings.Title(name)}
	case reflect.Ptr:
		if typ.Elem().Kind() != reflect.Struct {
			panic(fmt.Errorf("only struct pointer is supported, but got %v", typ))
		}
		return getJSONType(typ.Elem(), name)
	case reflect.Slice, reflect.Array:
		return JSONType{Kind: "array", Type: typ.Elem().Name()}
	default:
		panic(fmt.Errorf("unsupported type %s", typ.Kind()))
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
	success := schema.SuccessResponse(name, statusCode, body)
	if success.Body != nil {
		AddDefinition(defs, name+"Response", reflect.ValueOf(success.Body))
	}
	failures := schema.FailureResponses(name)
	for _, failure := range failures {
		AddDefinition(defs, name+"ResponseError"+strconv.Itoa(failure.StatusCode), reflect.ValueOf(failure.Body))
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
