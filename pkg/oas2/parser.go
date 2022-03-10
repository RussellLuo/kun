package oas2

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/RussellLuo/kun/gen/http/parser"
	"github.com/RussellLuo/kun/pkg/caseconv"
)

type JSONType struct {
	Kind        string
	Type        string
	Format      string
	Description string
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

type Parser struct {
	defs     map[string]Definition
	registry *GoTypeRegistry
}

func NewParser() *Parser {
	return &Parser{
		defs:     make(map[string]Definition),
		registry: NewGoTypeRegistry(),
	}
}

func (p *Parser) Definitions() map[string]Definition {
	return p.defs
}

func (p *Parser) AddDefinition(name string, value reflect.Value, embedded bool) {
	if _, ok := p.defs[name]; ok {
		// Ignore duplicated definitions implicitly.
		return
	}

	switch value.Kind() {
	case reflect.Struct:
		p.addStructDefinition(name, value, embedded)

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
			keyValue := p.addSubDefinition(keyString, value.MapIndex(key), false)

			properties = append(properties, Property{
				Name: keyString,
				Type: p.getJSONType(keyValue.Type(), caseconv.ToUpperCamelCase(keyString), ""),
			})
		}

		p.defs[name] = Definition{
			Type:                 "object",
			ItemTypeOrProperties: properties,
		}

	case reflect.Slice, reflect.Array:
		p.addArrayDefinition(name, value, false)

	case reflect.Ptr:
		elemType := value.Type().Elem()
		elem := reflect.New(elemType).Elem()
		p.AddDefinition(name, elem, embedded) // Always use the input name

	default:
		panic(fmt.Errorf("unsupported type %s", value.Kind()))
	}
}

func (p *Parser) addStructDefinition(name string, value reflect.Value, embedded bool) (properties []Property) {
	if isTime(value) {
		// Ignore this struct if it is a time value (of type `time.Time`).
		return
	}

	structType := value.Type()

	if !p.registry.Register(structType, name) {
		return p.registry.Properties(structType)
	}
	defer func() {
		p.registry.SetProperties(structType, properties)
	}()

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

		var fieldValueType reflect.Type

		structField := &parser.StructField{
			Name:      field.Name,
			//CamelCase: false,
			Type:      field.Type.Name(),
			Tag:       field.Tag,
		}
		if err := structField.Parse(); err != nil {
			panic(err)
		}

		if len(structField.Params) > 1 {
			// XXX: Add support for Argument aggregation.
			panic(fmt.Errorf("argument aggregation is unsupported in OAS now"))
		}
		structFieldParam := structField.Params[0]

		if structFieldParam.Type != structField.Type {
			// Use the user-specified type (a basic type) if any.
			var err error
			if fieldValueType, err = getReflectType(structFieldParam.Type); err != nil {
				panic(err)
			}
		} else {
			// Use the raw type of this struct field.
			fieldValue := p.addSubDefinition(fieldName, value.Field(i), field.Anonymous)
			fieldValueType = fieldValue.Type()
		}

		if field.Anonymous {
			// If this is an embedded field, promote the sub-properties of this field.

			var subProperties []Property

			ft := field.Type
			switch k := ft.Kind(); {
			case k == reflect.Struct:
				v := value.Field(i)
				subProperties = p.addStructDefinition("", v, field.Anonymous)
			case k == reflect.Ptr && ft.Elem().Kind() == reflect.Struct:
				v := reflect.New(ft.Elem()).Elem()
				subProperties = p.addStructDefinition("", v, field.Anonymous)
			}

			properties = append(properties, subProperties...)
		} else {
			// Otherwise, append this field as a property.
			properties = append(properties, Property{
				Name: fieldName,
				Type: p.getJSONType(fieldValueType, caseconv.ToUpperCamelCase(fieldName), structFieldParam.Description),
			})
		}
	}

	// Only add non-embedded struct into definitions.
	if !embedded {
		p.defs[name] = Definition{
			Type:                 "object",
			ItemTypeOrProperties: properties,
		}
	}

	return
}

func (p *Parser) addSubDefinition(name string, value reflect.Value, embedded bool) reflect.Value {
	typeName := value.Type().Name()
	if typeName == "" {
		typeName = caseconv.ToUpperCamelCase(name)
	}

	switch value.Kind() {
	case reflect.Struct:
		// We only need to call AddDefinition if this is a non-embedded struct.
		// Otherwise, another call to addStructDefinition will be triggered
		// instead within addStructDefinition.
		if !embedded {
			p.AddDefinition(typeName, value, embedded)
		}
	case reflect.Map:
		p.AddDefinition(typeName, value, false)
	case reflect.Slice, reflect.Array:
		p.addArrayDefinition(typeName, value, true)
	case reflect.Ptr:
		elemType := value.Type().Elem()
		elemName := elemType.Name()
		elem := reflect.New(elemType).Elem()
		if !isBasicKind(elem.Kind()) {
			// This is a pointer to a non-basic type, add more possible definitions.
			p.AddDefinition(elemName, elem, embedded)
		}
	case reflect.Interface:
		value = p.addSubDefinition(typeName, value.Elem(), embedded)
	}

	return value
}

func (p *Parser) addArrayDefinition(name string, value reflect.Value, inner bool) {
	elemType := value.Type().Elem()
	k := elemType.Kind()

	if isBasicKind(k) {
		if !inner {
			p.defs[name] = Definition{
				Type:                 "array",
				ItemTypeOrProperties: p.getJSONType(elemType, elemType.Name(), ""),
			}
		}
		return
	}

	switch k {
	case reflect.Struct, reflect.Map:
		elem := reflect.New(elemType).Elem()
		p.AddDefinition(p.getArrayElemTypeName(elemType, name), elem, false)
	case reflect.Ptr:
		elemType = elemType.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		elem := reflect.New(elemType).Elem()
		p.AddDefinition(p.getArrayElemTypeName(elemType, name), elem, false)
	case reflect.Slice, reflect.Array:
		elem := reflect.New(elemType).Elem()
		p.addArrayDefinition(p.getArrayElemTypeName(elemType, name), elem, inner)
	default:
		panic(fmt.Errorf("only struct slice or array is supported, but got %v", elemType.String()))
	}

	if !inner {
		p.defs[name] = Definition{
			Type:                 "array",
			ItemTypeOrProperties: p.getArrayElemTypeName(elemType, name),
		}
	}
}

func (p *Parser) getArrayElemTypeName(elemType reflect.Type, arrayTypeName string) string {
	return p.getTypeName(elemType, arrayTypeName+"ArrayItem")
}

func (p *Parser) getJSONType(typ reflect.Type, name, description string) JSONType {
	switch typ.Kind() {
	case reflect.Bool:
		return JSONType{Kind: "basic", Type: "boolean", Description: description}
	case reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return JSONType{Kind: "basic", Type: "integer", Format: "int32", Description: description}
	case reflect.Int, reflect.Int64,
		reflect.Uint, reflect.Uint64, reflect.Uintptr:
		return JSONType{Kind: "basic", Type: "integer", Format: "int64", Description: description}
	case reflect.Float32:
		return JSONType{Kind: "basic", Type: "number", Format: "float", Description: description}
	case reflect.Float64:
		return JSONType{Kind: "basic", Type: "number", Format: "double", Description: description}
	case reflect.String:
		return JSONType{Kind: "basic", Type: "string", Description: description}
	case reflect.Struct:
		if isTime(reflect.New(typ).Elem()) {
			// A time value is also a struct in Go, but it is represented as a string in OAS.
			return JSONType{Kind: "basic", Type: "string", Format: "date-time", Description: description}
		}
		return JSONType{Kind: "object", Type: p.getTypeName(typ, name), Description: description}
	case reflect.Map:
		return JSONType{Kind: "object", Type: name, Description: description}
	case reflect.Ptr:
		// Dereference the pointer and get its element type.
		return p.getJSONType(typ.Elem(), name, description)
	case reflect.Slice, reflect.Array:
		elemType := typ.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		return JSONType{Kind: "array", Type: p.getArrayElemTypeName(elemType, name), Description: description}
	default:
		panic(fmt.Errorf("unsupported type %s", typ.Kind()))
	}
}

func (p *Parser) getTypeName(t reflect.Type, defaultName string) string {
	name := p.registry.Name(t)
	if name != "" {
		return name
	}

	if t.Name() != "" {
		return t.Name()
	}
	return defaultName
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

func isTime(v reflect.Value) bool {
	switch v.Interface().(type) {
	case *time.Time, time.Time:
		return true
	default:
		return false
	}
}

func getReflectType(typ string) (reflect.Type, error) {
	var v interface{}
	switch typ {
	case "bool":
		v = false
	case "string":
		v = ""
	case "int":
		v = int(0)
	case "int8":
		v = int8(0)
	case "int16":
		v = int16(0)
	case "int32":
		v = int32(0)
	case "int64":
		v = int64(0)
	case "uint":
		v = uint(0)
	case "uint16":
		v = uint16(0)
	case "uint32":
		v = uint32(0)
	case "uint64":
		v = uint64(0)
	case "float32":
		v = float32(0)
	case "float64":
		v = float64(0)
	case "time":
		v = time.Time{}
	default:
		return nil, fmt.Errorf("invalid basic type name: %s", typ)
	}
	return reflect.ValueOf(v).Type(), nil
}
