package parser

import (
	"fmt"
	"go/types"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/RussellLuo/kun/gen/http/parser/annotation"
	"github.com/RussellLuo/kun/gen/http/spec"
	utilannotation "github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/docutil"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

const (
	OptionNoRequestBody = "-"
)

var (
	rePathVarName = regexp.MustCompile(`{(\w+)}`)
)

func Parse(data *ifacetool.Data, snakeCase bool) (*spec.Specification, docutil.Transport, error) {
	anno, err := annotation.ParseInterfaceAnnotation(data.InterfaceDoc)
	if err != nil {
		return nil, 0, err
	}

	s := &spec.Specification{
		Metadata: anno.Metadata,
	}

	var (
		transport docutil.Transport
		opBuilder = &OpBuilder{
			snakeCase: snakeCase,
			aliases:   anno.Aliases,
		}
	)

	for _, m := range data.Methods {
		doc := docutil.Doc(m.Doc).JoinComments()
		m.Doc = doc // Replace the original doc with joined doc.

		t := doc.Transport()
		if t == 0 {
			// Empty transport indicates that there are no annotations.
			continue
		}
		transport |= t

		if !t.Has(docutil.TransportHTTP) {
			// Add operations for generating endpoint code for gRPC.
			op := spec.NewOperation(m.Name, m.Name, annotation.GetDescriptionFromDoc(m.Doc))
			for _, arg := range m.Params {
				op.Request.Bind(arg, opBuilder.buildParams(arg, nil))
			}
			s.Operations = append(s.Operations, op)
			continue
		}

		ops, err := opBuilder.Build(m)
		if err != nil {
			return nil, 0, err
		}

		s.Operations = append(s.Operations, ops...)
	}

	return s, transport, nil
}

type OpBuilder struct {
	snakeCase bool
	aliases   annotation.Aliases
}

func (b *OpBuilder) Build(method *ifacetool.Method) ([]*spec.Operation, error) {
	anno, err := annotation.ParseMethodAnnotation(method, b.aliases)
	if err != nil {
		return nil, err
	}

	if len(anno.Ops) == 0 {
		return nil, fmt.Errorf("method %s has no comment directive %s", method.Name, utilannotation.DirectiveHTTPOp)
	}

	var allPathVarNames PathVarNames
	for _, annoOp := range anno.Ops {
		allPathVarNames.Add(extractPathVarNames(annoOp.Pattern))
	}

	var ops []*spec.Operation
	for i, annoOp := range anno.Ops {
		name := method.Name
		if i > 0 {
			// Append a suffix to the name from the second operation, if any,
			// to differentiate one operation from another.
			name += strconv.Itoa(i)
		}
		op := spec.NewOperation(name, method.Name, annotation.GetDescriptionFromDoc(method.Doc))
		{
			// Set the HTTP method and the URI pattern.
			op.Method, op.Pattern = annoOp.Method, annoOp.Pattern

			// Set the request body field.
			//
			// We must do this before setting request parameters, in case no-request-body is specified.
			b.setBodyField(op.Request, anno.Body)

			// Set the request parameters.
			//
			// Note that the way to handle path parameters here:
			//   - First, we set all path parameters collected from all patterns in anno.Ops;
			//   - Then we remove the path parameters which are not defined in op's own pattern.
			//
			// The reason for doing so, is to properly handle the cases where path parameters
			// are specified explicitly, either in annotations or in struct tags. For example:
			//
			// ```go
			// type GetMessageRequest struct {
			//     userID    string `kun:"in=path name=userID"`
			//     messageID string `kun:"in=path name=messageID"`
			// }
			//
			// type Service interface {
			//     //kun:op GET /messages/{messageID}
			//     //kun:op GET /users/{userID}/messages/{messageID}
			//     //kun:param req
			//     GetMessage(ctx context.Context, req GetMessageRequest) (text string, err error)
			// }
			// ```
			//
			// If we don't pass in all path parameters, for `//kun:op GET /messages/{messageID}`,
			// the binding specified by `kun:"in=path name=userID"` will trigger an error since
			// the associated path parameter `userID` is not defined in the corresponding pattern.
			if err := b.setParams(op.Request, method, anno.Params, allPathVarNames.Squash()); err != nil {
				return nil, err
			}
			op.Request.Bindings = removePathParamsNotItsOwn(op.Request.Bindings, allPathVarNames.Get(i))

			// Manipulate the request body.
			if err := b.manipulateBody(op.Request, anno.Body); err != nil {
				return nil, err
			}

			// Set the success response.
			if anno.Success != nil {
				op.SuccessResponse = anno.Success
			}

			// Set the OAS tags.
			op.Tags = anno.Tags
		}

		ops = append(ops, op)
	}

	return ops, nil
}

func (b *OpBuilder) setBodyField(req *spec.Request, body *annotation.Body) {
	if body == nil {
		return
	}

	if body.Field != "" {
		req.BodyField = body.Field
	}
}

func (b *OpBuilder) manipulateBody(req *spec.Request, body *annotation.Body) error {
	if body == nil {
		return nil
	}

	if req.BodyField != "" {
		if len(body.Manipulations) > 0 {
			if req.BodyField == OptionNoRequestBody {
				return fmt.Errorf("useless manipulations in %s since there is no request body", utilannotation.DirectiveHTTPBody)
			}
			return fmt.Errorf("useless manipulations in %s since the request body has been mapped to argument %q", utilannotation.DirectiveHTTPBody, req.BodyField)
		}
		return nil
	}

	for _, binding := range req.Bindings {
		m, ok := body.Manipulations[binding.Arg.Name]
		if !ok {
			continue
		}

		if binding.In() != spec.InBody {
			return fmt.Errorf("argument %q manipulated in %s is not located in body", binding.Arg.Name, utilannotation.DirectiveHTTPBody)
		}

		binding.SetName(m.Name)
		// TODO: Modify the endpoint generator to add annotations (about OAS type and description) in request struct.
		binding.SetType(m.Type)
		binding.SetDescription(m.Description)
	}

	return nil
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

	// TODO: ensure that all path parameters — specified explicitly — are defined in the pattern.

	// Automatically bind path parameters defined in the pattern.
	for _, name := range pathVarNames {
		// If name is already bound to a path parameter by //kun:param or
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
	if req.BodyField == OptionNoRequestBody {
		for _, binding := range req.Bindings {
			if !binding.IsManual() {
				// Rebind binding.Arg as HTTP request parameters.
				annoParams, err := b.inferAnnotationParams(method.Name, binding.Arg)
				if err != nil {
					return err
				}
				binding.Params = b.buildParams(binding.Arg, annoParams)
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

// inferAnnotationParams extracts parameters by parsing annotations from struct tags,
// which are only used for HTTP request parameters.
func (b *OpBuilder) inferAnnotationParams(methodName string, arg *ifacetool.Param) ([]*spec.Parameter, error) {
	newParamWithType := func(typ string) *spec.Parameter {
		return &spec.Parameter{
			// Method arguments specified in //kun:param are mapped to the query by default.
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

		for i := 0; i < t.NumFields(); i++ {
			var typeName string
			switch ft := t.Field(i).Type().Underlying().(type) {
			case *types.Basic:
				typeName = ft.Name()
			case *types.Slice:
				et, ok := ft.Elem().(*types.Basic)
				if !ok {
					// XXX: We should allow this type if it's used in Argument aggregation.
					return nil, newErrorStruct(arg.Name, t.Field(i).Name(), ft)
				}
				typeName = "[]" + et.Name()
			default:
				// XXX: We should allow this type if it's used in Argument aggregation.
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
				continue
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

	case *types.Interface:
		// We assume this is `context.Context`, just ignore it.
		return nil, nil

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

type PathVarNames [][]string

func (pvn *PathVarNames) Add(names []string) {
	*pvn = append(*pvn, names)
}

func (pvn *PathVarNames) Get(i int) []string {
	if i >= len(*pvn) {
		return nil
	}
	return (*pvn)[i]
}

// Squash de-duplicates all the nested names and squash them into a flat list.
func (pvn *PathVarNames) Squash() []string {
	var flat []string
	for _, names := range *pvn {
		for _, n := range names {
			if !sliceContains(flat, n) {
				flat = append(flat, n)
			}
		}
	}
	return flat
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

func sliceContains(slice []string, target string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

// removePathParamsNotItsOwn removes all path parameters, which are not
// contained in pathVarNames.
func removePathParamsNotItsOwn(bindings []*spec.Binding, pathVarNames []string) []*spec.Binding {
	var result []*spec.Binding

	for _, b := range bindings {
		switch {
		case b.IsAggregate():
			var params []*spec.Parameter
			for _, p := range b.Params {
				if p.In != spec.InPath || sliceContains(pathVarNames, p.Name) {
					params = append(params, p)
				}
			}
			if len(params) > 0 {
				result = append(result, &spec.Binding{
					Arg:    b.Arg,
					Params: params,
				})
			}
		case b.In() != spec.InPath || sliceContains(pathVarNames, b.Name()):
			result = append(result, b)
		}
	}

	return result
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
	params, err := annotation.ParseParamOptions(f.Name, f.Tag.Get(utilannotation.Name))
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
