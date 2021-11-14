package parser

import (
	"fmt"
	"go/types"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/http/parser/annotation"
	"github.com/RussellLuo/kok/gen/http/spec"
	"github.com/RussellLuo/kok/pkg/caseconv"
	"github.com/RussellLuo/kok/pkg/ifacetool"
)

type Transport int

const (
	OptionNoBody = "-"

	TransportHTTP Transport = 0b0001
	TransportGRPC Transport = 0b0010
	TransportAll  Transport = 0b0011

	tagName = "kok"
)

var (
	rePathVarName = regexp.MustCompile(`{(\w+)}`)
)

func Parse(data *ifacetool.Data, snakeCase bool) (*spec.Specification, []Transport, error) {
	anno, err := annotation.ParseInterfaceAnnotation(data.InterfaceDoc)
	if err != nil {
		return nil, nil, err
	}

	s := &spec.Specification{
		Metadata: anno.Metadata,
	}

	var (
		transports []Transport
		opBuilder  = &OpBuilder{
			snakeCase: snakeCase,
			aliases:   anno.Aliases,
		}
	)

	for _, m := range data.Methods {
		transport := getTransportPerKokAnnotations(m.Doc)
		if transport == 0 {
			// Empty transport indicates that there are no kok annotations.
			continue
		}
		transports = append(transports, transport)

		op, err := opBuilder.Build(m)
		if err != nil {
			return nil, nil, err
		}

		s.Operations = append(s.Operations, op)
	}

	return s, transports, nil
}

type OpBuilder struct {
	snakeCase bool
	aliases   annotation.Aliases
}

func (b *OpBuilder) Build(method *ifacetool.Method) (*spec.Operation, error) {
	op := &spec.Operation{
		Name:        method.Name,
		Description: annotation.GetDescriptionFromDoc(method.Doc),
	}

	anno, err := annotation.ParseMethodAnnotation(method)
	if err != nil {
		return nil, err
	}

	// Set method and pattern.
	op.Method, op.Pattern = anno.Op.Method, anno.Op.Pattern
	if op.Method == "" && op.Pattern == "" {
		return nil, fmt.Errorf("method %s has no comment directive //kok:op", method.Name)
	}

	// Set request parameters.
	op.Request = new(spec.Request)
	pathVarNames := extractPathVarNames(op.Pattern)
	if err := b.setParams(op.Request, method, anno.Params, pathVarNames); err != nil {
		return nil, err
	}

	// Set request body.
	if anno.Body != nil {
		op.Request.BodyField = anno.Body.Field
	}

	// Set success response.
	op.Resp(http.StatusOK, spec.MediaTypeJSON, nil)
	if anno.Success != nil {
		op.SuccessResponse = anno.Success
	}

	// Set OAS tags.
	op.Tags = anno.Tags

	return op, nil
}

func (b *OpBuilder) setParams(req *spec.Request, method *ifacetool.Method, params map[string]*annotation.Param, pathVarNames []string) error {
	for _, arg := range method.Params {
		param, ok := params[arg.Name]
		if !ok {
			req.Bind(arg, b.buildParams(arg, nil))
			continue
		}

		// Remove this entry to check unmatched annotations later.
		delete(params, arg.Name)

		if len(param.Params) > 0 {
			if !isBasic(arg) {
				return fmt.Errorf("cannot define extra parameters for non-basic argument %q", arg.Name)
			}
			req.Bind(arg, b.buildParams(arg, param.Params))
			continue
		}

		annoParams, err := b.inferAnnotationParams(method.Name, arg)
		if err != nil {
			return err
		}
		req.Bind(arg, b.buildParams(arg, annoParams))
	}

	for name, p := range params {
		// Remain some unmatched entries.

		if !strings.HasPrefix(name, "__") {
			return fmt.Errorf("no argument %q declared in the method %s", name, method.Name)
		}

		// This is a blank identifier.
		typ := types.Typ[types.String] // Defaults to string.
		arg := &ifacetool.Param{
			Name:       name,
			TypeString: typ.Name(),
			Type:       typ,
		}
		req.Bind(arg, b.buildParams(arg, p.Params))
	}

	// Add path parameters according to the path pattern.
	for _, name := range pathVarNames {
		// If name is already bound to a path parameter by //kok:param or
		// by struct tags, do not reset it.
		if isAlreadyPathParam(name, req.Bindings) {
			continue
		}

		// Get the method argument name by always converting the path variable
		// name to lowerCamelCase (the naming convention in Go).
		//
		// Known issues:
		// - "xx_id" will be converted to "xxId" (not the conventional "xxID").
		argName := caseconv.ToLowerCamelCase(name)

		binding := req.GetBinding(argName)
		if binding == nil {
			return fmt.Errorf("cannot bind path parameter %q: no argument %q declared in the method %s", name, argName, method.Name)
		}

		if binding.IsManual() {
			return fmt.Errorf("cannot bind path parameter %q: argument %q has been mannually bound to a %s parameter named %q", name, argName, binding.In(), binding.Name())
		}

		// Rebind it to a path parameter named name.
		binding.SetIn(spec.InPath)
		binding.SetName(name)
	}

	// Add possible query parameters if no-request-body is specified.
	if req.BodyField == OptionNoBody {
		for _, binding := range req.Bindings {
			if !binding.IsManual() {
				binding.SetIn(spec.InQuery)
			}
		}
	}

	return nil
}

func (b *OpBuilder) buildParams(arg *ifacetool.Param, annoParams []*spec.Parameter) []*spec.Parameter {
	defaultParam := &spec.Parameter{
		In:   spec.InBody, // Method arguments are bound to the body by default.
		Name: b.defaultName(arg.Name),
		Type: arg.TypeString,
	}

	if len(annoParams) == 0 {
		// No annotation parameter is specified, use the default one.
		return []*spec.Parameter{defaultParam}
	}

	// Complete the properties of annoParams by using values from defaultParam.
	for _, p := range annoParams {
		if p.Name == "" {
			p.Name = defaultParam.Name
		}
		if p.Type == "" {
			p.Type = defaultParam.Type
		}
	}
	return annoParams
}

// inferAnnotationParams extracts parameters by parsing annotations from struct tags.
func (b *OpBuilder) inferAnnotationParams(methodName string, arg *ifacetool.Param) ([]*spec.Parameter, error) {
	newParamWithType := func(typ string) *spec.Parameter {
		return &spec.Parameter{
			// Method arguments specified in //kok:param are mapped to the query by default.
			In:   spec.InQuery,
			Name: b.defaultName(arg.Name),
			Type: typ,
		}
	}

	newError := func(name string, typ types.Type) error {
		return fmt.Errorf("parameter cannot be mapped to argument %q (of type %v) in method %s", name, typ, methodName)
	}

	newErrorStruct := func(paramName, fieldName string, fieldType types.Type) error {
		return fmt.Errorf("parameter cannot be mapped to struct field %q (of type %v) from argument %q in method %s", fieldName, fieldType, paramName, methodName)
	}

	var params []*spec.Parameter

	switch t := arg.Type.Underlying().(type) {
	case *types.Basic:
		params = append(params, newParamWithType(t.Name()))

	case *types.Slice:
		et, ok := t.Elem().(*types.Basic)
		if !ok {
			return nil, newError(arg.Name, t)
		}
		params = append(params, newParamWithType("[]"+et.Name()))

	case *types.Struct:
		/*if a.In != "" {
			fmt.Printf("WARNING: manually specified `in:%s` is ignored for struct argument `%s` in method %s\n", a.In, a.ArgName, p.methodName)
		}*/

	NextField:
		for i := 0; i < t.NumFields(); i++ {
			var typeName string
			switch ft := t.Field(i).Type().(type) {
			case *types.Basic:
				typeName = ft.Name()
			case *types.Slice:
				et, ok := ft.Elem().(*types.Basic)
				if !ok {
					return nil, newErrorStruct(arg.Name, t.Field(i).Name(), ft)
				}
				typeName = "[]" + et.Name()
			default:
				return nil, newErrorStruct(arg.Name, t.Field(i).Name(), ft)
			}

			field := &StructField{
				Name:      t.Field(i).Name(),
				CamelCase: !b.snakeCase,
				Type:      typeName,
				Tag:       reflect.StructTag(t.Tag(i)),
			}
			if err := field.Parse(); err != nil {
				return nil, err
			}

			if field.Omitted {
				// Omit this field.
				continue NextField
			}

			params = append(params, field.Params...)
		}

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
		return nil, newError(arg.Name, t)
	}

	return params, nil
}

func (b *OpBuilder) defaultName(name string) string {
	if b.snakeCase {
		return caseconv.ToSnakeCase(name)
	}
	return caseconv.ToLowerCamelCase(name)
}

// isBasic returns whether the parameter is of basic type, or of
// slice type (whose element is of basic type).
func isBasic(arg *ifacetool.Param) bool {
	switch t := arg.Type.Underlying().(type) {
	case *types.Basic:
		return true

	case *types.Slice:
		_, ok := t.Elem().(*types.Basic)
		return ok
	}
	return false
}

func getTransportPerKokAnnotations(doc []string) (t Transport) {
	for _, comment := range doc {
		if annotation.IsKokGRPCAnnotation(comment) {
			t = t | TransportGRPC
		} else if annotation.IsKokAnnotation(comment) {
			t = t | TransportHTTP
		}
	}
	return t
}

func extractPathVarNames(pattern string) (names []string) {
	result := rePathVarName.FindAllStringSubmatch(pattern, -1)
	if len(result) == 0 {
		return
	}

	for _, s := range result {
		names = append(names, s[1])
	}
	return
}

func isAlreadyPathParam(name string, bindings []*spec.Binding) bool {
	for _, b := range bindings {
		for _, p := range b.Params {
			if p.In == spec.InPath && p.Name == name {
				return true
			}
		}
	}
	return false
}

type StructField struct {
	Name      string
	CamelCase bool
	Type      string
	Tag       reflect.StructTag

	Omitted bool              // Whether to omit this field.
	Params  []*spec.Parameter // The associated annotation parameters.
}

func (f *StructField) Parse() error {
	params, err := annotation.ParseParamParameters(f.Name, f.Tag.Get(tagName))
	if err != nil {
		return err
	}

	for _, p := range params {
		if p.Name == "-" {
			f.Omitted = true
			return nil
		}

		if p.Name == "" {
			p.Name = caseconv.ToSnakeCase(f.Name)
			if f.CamelCase {
				p.Name = caseconv.ToLowerCamelCase(f.Name)
			}
		}
		if p.Type == "" {
			p.Type = f.Type
		}
	}

	f.Params = params
	return nil
}
