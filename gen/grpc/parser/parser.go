package parser

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/RussellLuo/kok/gen/util/reflector"
	"github.com/RussellLuo/kok/pkg/caseconv"
	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
)

var (
	// Go Type -> .proto Type
	// see https://developers.google.com/protocol-buffers/docs/proto3#scalar
	scalarTypes = map[string]string{
		"float64": "double",
		"float32": "float",
		"int32":   "int32",  // sint32, sfixed32
		"int64":   "int64",  // sint64, sfixed64
		"uint32":  "uint32", // fixed32
		"uint64":  "uint64", // fixed64
		"bool":    "bool",
		"string":  "string",
		"[]byte":  "bytes",
	}
)

type Service struct {
	Name         string
	RPCs         []*RPC
	Descriptions []string
}

type RPC struct {
	Name         string
	Request      *Message
	Response     *Message
	Descriptions []string
}

type Message struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name string
	Type *Type
	Num  int
}

type Type struct {
	Name     string
	Repeated bool
	Fields   []*Field
}

// Squash does a pre-order walk of t and returns all the composite types
// (including itself) as a flat list.
func (t *Type) Squash() (types []*Type) {
	if len(t.Fields) > 0 {
		types = append(types, t)
	}
	for _, f := range t.Fields {
		if len(f.Type.Fields) > 0 {
			types = append(types, f.Type)
			types = append(types, f.Type.Squash()...)
		}
	}
	return
}

func Parse(result *reflector.Result, doc *reflector.InterfaceDoc) (*Service, error) {
	s := &Service{
		Name:         result.Interface.Name,
		Descriptions: getDescriptionsFromDoc(doc.Doc),
	}

	for _, m := range result.Interface.Methods {
		comments, ok := doc.MethodDocs[m.Name]
		if !ok || !hasKokGRPCAnnotation(comments) {
			continue
		}

		reqFields, err := parse(m.Params)
		if err != nil {
			return nil, err
		}

		respFields, err := parse(m.Returns)
		if err != nil {
			return nil, err
		}

		s.RPCs = append(s.RPCs, &RPC{
			Name:         m.Name,
			Descriptions: getDescriptionsFromDoc(comments),
			Request: &Message{
				Name:   m.Name + "Request",
				Fields: reqFields,
			},
			Response: &Message{
				Name:   m.Name + "Response",
				Fields: respFields,
			},
		})
	}

	return s, nil
}

func parse(params []*reflector.Param) ([]*Field, error) {
	var fields []*Field
	var i int
	for _, p := range params {
		if p.Type == "context.Context" || p.Type == "error" {
			continue
		}

		typ, err := parseType(p.Name, p.RawType)
		if err != nil {
			return nil, err
		}

		i++
		fields = append(fields, &Field{
			Name: caseconv.ToSnakeCase(p.Name),
			Type: typ,
			Num:  i,
		})
	}
	return fields, nil
}

func parseType(name string, typ types.Type) (*Type, error) {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		return parseBasicType(t), nil

	case *types.Slice:
		st, err := parseSliceType(name, t)
		if err != nil {
			return nil, err
		}
		return st, nil

	case *types.Struct:
		st, err := parseStructType(name, typ, t)
		if err != nil {
			return nil, err
		}
		return st, nil

		//case *types.Pointer:
		//	// Dereference the pointer to parse the element type.
		//	nts, err := p.parseTypes(name, t.Elem())
		//	if err != nil {
		//		return nil, err
		//	}
		//	for n, t := range nts {
		//		nameTypes[n] = t
		//	}
	default:
		return nil, fmt.Errorf("unsupported %T", t)
	}
}

func parseBasicType(t *types.Basic) *Type {
	return &Type{Name: scalarTypes[t.Name()]}
}

func parseSliceType(name string, t *types.Slice) (*Type, error) {
	typ, err := parseType(name, t.Elem())
	if err != nil {
		return nil, err
	}

	if typ.Name == "byte" {
		// Go: []byte => proto: bytes
		return &Type{Name: "bytes"}, nil
	}

	return &Type{Name: typ.Name, Repeated: true}, nil
}

func parseStructType(name string, typ types.Type, t *types.Struct) (*Type, error) {
	// Try to get the actual type name if typ is a named type.
	named, ok := typ.(*types.Named)
	if ok {
		name = named.Obj().Name()
	}

	var fields []*Field
	for i := 0; i < t.NumFields(); i++ {
		fieldName := getFieldName(t, i)
		if fieldName == "" {
			// Empty name indicates omitting.
			return nil, nil
		}

		fieldType, err := getFieldType(t, i, fieldName)
		if err != nil {
			return nil, err
		}

		fields = append(fields, &Field{
			Name: caseconv.ToSnakeCase(fieldName),
			Type: fieldType,
			Num:  i + 1,
		})
	}
	return &Type{Name: caseconv.ToUpperCamelCase(name), Fields: fields}, nil
}

func getFieldName(t *types.Struct, i int) string {
	field := reflect.StructField{
		Tag:  reflect.StructTag(t.Tag(i)),
		Name: t.Field(i).Name(),
	}

	kokField := httpcodec.GetKokField(field)
	if kokField.Omitted {
		return ""
	}

	parts := strings.SplitN(kokField.Name, ".", 2)
	return parts[1]
}

func getFieldType(t *types.Struct, i int, fieldName string) (*Type, error) {
	typ := t.Field(i).Type()

	switch ft := typ.Underlying().(type) {
	case *types.Basic:
		return parseBasicType(ft), nil

	case *types.Slice:
		st, err := parseSliceType(fieldName, ft)
		if err != nil {
			return nil, err
		}
		return st, nil

	case *types.Struct:
		st, err := parseStructType(fieldName, typ, ft)
		if err != nil {
			return nil, err
		}
		return st, nil

	default:
		return nil, fmt.Errorf("unsupported %T", ft)
	}
}

func getDescriptionsFromDoc(doc []string) (comments []string) {
	for _, comment := range doc {
		if !isKokAnnotation(comment, "@kok") && !isKokAnnotation(comment, "@kok(grpc)") {
			comments = append(comments, comment)
		}
	}
	return
}

func hasKokGRPCAnnotation(comments []string) bool {
	for _, comment := range comments {
		if isKokAnnotation(comment, "@kok(grpc)") {
			return true
		}
	}
	return false
}

func isKokAnnotation(comment, anno string) bool {
	content := strings.TrimPrefix(comment, "//")
	trimmed := strings.TrimSpace(content)
	return strings.HasPrefix(trimmed, anno)
}
