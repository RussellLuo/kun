package parser

import (
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"strings"

	"github.com/RussellLuo/kun/gen/http/parser"
	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
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
		"int":     "int64",
		"bool":    "bool",
		"string":  "string",
		"[]byte":  "bytes",
	}

	reGRPC = regexp.MustCompile(`^` + annotation.DirectiveGRPC.String() + `(.*)$`)
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
	Repeated bool     // true for slice: []Type
	MapKey   string   // non-empty for map: map[key]Type
	Fields   []*Field // non-empty for struct
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

func Parse(data *ifacetool.Data) (*Service, error) {
	s := &Service{
		Name:         data.InterfaceName,
		Descriptions: getDescriptionsFromDoc(data.InterfaceDoc),
	}

	for _, m := range data.Methods {
		if len(m.Doc) == 0 || !hasGRPCAnnotation(m.Doc) {
			continue
		}

		rpcFields, err := parseRPCFields(m)
		if err != nil {
			return nil, err
		}

		s.RPCs = append(s.RPCs, &RPC{
			Name:         m.Name,
			Descriptions: getDescriptionsFromDoc(m.Doc),
			Request: &Message{
				Name:   m.Name + "Request",
				Fields: rpcFields.Request,
			},
			Response: &Message{
				Name:   m.Name + "Response",
				Fields: rpcFields.Response,
			},
		})
	}

	return s, nil
}

func parse(params []*ifacetool.Param) ([]*Field, error) {
	var fields []*Field
	var i int
	for _, p := range params {
		if p.TypeString == "context.Context" || p.TypeString == "error" {
			continue
		}

		typ, err := parseType(p.Name, p.Type)
		if err != nil {
			return nil, err
		}

		i++
		fields = append(fields, &Field{
			Name: p.Name,
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

	case *types.Map:
		k := t.Key()
		bt, ok := k.(*types.Basic)
		if !ok {
			return nil, fmt.Errorf("unsupported map key %T", k)
		}
		kt := parseBasicType(bt)

		// TODO: Add support for map[string]interface{}?
		vt, err := parseType("", t.Elem())
		if err != nil {
			return nil, err
		}

		return &Type{
			Name:   vt.Name,
			MapKey: kt.Name,   // type name of the map key
			Fields: vt.Fields, // possible fields from the map value.
		}, nil

	case *types.Struct:
		st, err := parseStructType(name, typ, t)
		if err != nil {
			return nil, err
		}
		return st, nil

	case *types.Pointer:
		// Dereference the pointer to parse the element type.
		et, err := parseType(name, t.Elem())
		if err != nil {
			return nil, err
		}
		return et, nil

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

		fieldType, err := parseType(fieldName, t.Field(i).Type())
		if err != nil {
			return nil, err
		}

		fields = append(fields, &Field{
			Name: fieldName,
			Type: fieldType,
			Num:  i + 1,
		})
	}
	return &Type{Name: caseconv.ToUpperCamelCase(name), Fields: fields}, nil
}

func getFieldName(t *types.Struct, i int) string {
	field := &parser.StructField{
		Name: t.Field(i).Name(),
		Tag:  reflect.StructTag(t.Tag(i)),
	}
	if err := field.Parse(); err != nil {
		return ""
	}

	if field.Omitted {
		return ""
	}

	return field.Name
}

func getDescriptionsFromDoc(doc []string) (comments []string) {
	for _, comment := range doc {
		if annotation.Directive(comment).Dialect() == annotation.DialectUnknown {
			comments = append(comments, comment)
		}
	}
	return
}

func hasGRPCAnnotation(doc []string) bool {
	for _, comment := range doc {
		if annotation.Directive(comment).Dialect() == annotation.DialectGRPC {
			return true
		}
	}
	return false
}

type rpcFields struct {
	Request  []*Field
	Response []*Field
}

func parseRPCFields(method *ifacetool.Method) (*rpcFields, error) {
	reqFields, err := parse(method.Params)
	if err != nil {
		return nil, err
	}

	respFields, err := parse(method.Returns)
	if err != nil {
		return nil, err
	}

	rpcFields := &rpcFields{
		Request:  reqFields,
		Response: respFields,
	}

	if err := rpcFields.manipulateByComments(method); err != nil {
		return nil, err
	}

	return rpcFields, nil
}

func (rf *rpcFields) manipulateByComments(method *ifacetool.Method) error {
	params := make(map[string]*ifacetool.Param)
	for _, p := range method.Params {
		params[p.Name] = p
	}

	returns := make(map[string]*ifacetool.Param)
	for _, p := range method.Returns {
		returns[p.Name] = p
	}

	for _, comment := range method.Doc {
		if annotation.Directive(comment).Dialect() != annotation.DialectGRPC {
			continue
		}

		result := reGRPC.FindStringSubmatch(comment)
		if len(result) != 2 {
			return fmt.Errorf("invalid %s directive: %s", annotation.DirectiveGRPC, comment)
		}
		value := strings.TrimSpace(result[1])
		if value == "" {
			continue
		}

		fields := strings.Fields(value)
		for _, f := range fields {
			parts := strings.Split(f, "=")
			if len(parts) != 2 {
				return fmt.Errorf(`%q does not match the expected format: <key>=<value>`, f)
			}
			k, v := parts[0], parts[1]

			switch k {
			case "request":
				p, ok := params[v]
				if !ok {
					return fmt.Errorf("no param `%s` declared in the method %s", v, method.Name)
				}
				if !isStructType(p.Type) {
					return fmt.Errorf("non-struct param `%s` in the method %s cannot be mapped to a gRPC request", v, method.Name)
				}

				structType, err := parseType(p.Name, p.Type)
				if err != nil {
					return err
				}
				rf.Request = structType.Fields

			case "response":
				p, ok := returns[v]
				if !ok {
					return fmt.Errorf("no result `%s` declared in the method %s", v, method.Name)
				}
				if !isStructType(p.Type) {
					return fmt.Errorf("non-struct result `%s` in the method %s cannot be mapped to a gRPC response", v, method.Name)
				}

				structType, err := parseType(p.Name, p.Type)
				if err != nil {
					return err
				}
				rf.Response = structType.Fields

			default:
				return fmt.Errorf(`unrecognized %s key "%s" in comment: %s`, annotation.Name, k, comment)
			}
		}
	}

	return nil
}

func isStructType(typ types.Type) bool {
	switch t := typ.Underlying().(type) {
	case *types.Struct:
		return true
	case *types.Pointer:
		return isStructType(t.Elem())
	default:
		return false
	}
}
