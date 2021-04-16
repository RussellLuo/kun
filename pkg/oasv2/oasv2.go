package oasv2

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
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
        {{- end}} {{/* if .Type.Format */}}

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
          {{- end}} {{/* if isBasic .Type.Type */}}
        {{- end}} {{/* if eq .Type.Kind "basic" */}}

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
		if isTime(value) {
			// Ignore this struct if it is a time value (of type `time.Time`).
			return
		}

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
		addArrayDefinition(defs, name, value, false)

	case reflect.Ptr:
		elemType := value.Type().Elem()
		elem := reflect.New(elemType).Elem()
		AddDefinition(defs, name, elem) // Always use the input name

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
	case reflect.Struct, reflect.Map:
		AddDefinition(defs, typeName, value)
	case reflect.Slice, reflect.Array:
		addArrayDefinition(defs, typeName, value, true)
	case reflect.Ptr:
		elemType := value.Type().Elem()
		elemName := elemType.Name()
		elem := reflect.New(elemType).Elem()
		if !isBasicKind(elem.Kind()) {
			// This is a pointer to a non-basic type, add more possible definitions.
			AddDefinition(defs, elemName, elem)
		}
	case reflect.Interface:
		value = addSubDefinition(defs, typeName, value.Elem())
	}

	return value
}

func addArrayDefinition(defs map[string]Definition, name string, value reflect.Value, inner bool) {
	elemType := value.Type().Elem()
	k := elemType.Kind()

	if isBasicKind(k) {
		if !inner {
			defs[name] = Definition{
				Type:                 "array",
				ItemTypeOrProperties: getJSONType(elemType, elemType.Name()),
			}
		}
		return
	}

	switch k {
	case reflect.Struct, reflect.Map:
		elem := reflect.New(elemType).Elem()
		AddDefinition(defs, elemType.Name(), elem)
	case reflect.Ptr:
		elemType = elemType.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		elem := reflect.New(elemType).Elem()
		AddDefinition(defs, elemType.Name(), elem)
	default:
		panic(fmt.Errorf("only struct slice or array is supported, but got %v", elemType.String()))
	}

	if !inner {
		defs[name] = Definition{
			Type:                 "array",
			ItemTypeOrProperties: elemType.Name(),
		}
	}
}

func isBasicKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
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
		if isTime(reflect.New(typ).Elem()) {
			// A time value is also a struct in Go, but it is represented as a string in OAS.
			return JSONType{Kind: "basic", Type: "string", Format: "date-time"}
		}
		return JSONType{Kind: "object", Type: typ.Name()}
	case reflect.Map:
		return JSONType{Kind: "object", Type: strings.Title(name)}
	case reflect.Ptr:
		// Dereference the pointer and get its element type.
		return getJSONType(typ.Elem(), name)
	case reflect.Slice, reflect.Array:
		elemType := typ.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		return JSONType{Kind: "array", Type: elemType.Name()}
	default:
		panic(fmt.Errorf("unsupported type %s", typ.Kind()))
	}
}

func isTime(v reflect.Value) bool {
	switch v.Interface().(type) {
	case *time.Time, time.Time:
		return true
	default:
		return false
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
