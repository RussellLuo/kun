package parser

import (
	"fmt"
	"go/types"
	"net/http"
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
)

var (
	rePathVarName = regexp.MustCompile(`{(\w+)}`)
)

func Parse(data *ifacetool.Data, snakeCase bool) (*spec.Specification, []Transport, error) {
	ifaceAnno, err := annotation.ParseInterfaceAnnotation(data.InterfaceDoc)
	if err != nil {
		return nil, nil, err
	}

	s := &spec.Specification{
		Metadata: ifaceAnno.Metadata,
	}

	var (
		transports []Transport
		opBuilder  = &OpBuilder{
			snakeCase: snakeCase,
			aliases:   ifaceAnno.Aliases,
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
	if err := b.setParams(op, method, anno.Params); err != nil {
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

func (b *OpBuilder) setParams(op *spec.Operation, method *ifacetool.Method, params map[string]*annotation.Param) error {
	for _, arg := range method.Params {
		param, ok := params[arg.Name]
		if !ok {
			op.Bind(arg, b.buildParams(arg, nil))
			continue
		}

		// Remove this entry to check unmatched annotations later.
		delete(params, arg.Name)

		if len(param.Params) > 0 {
			op.Bind(arg, b.buildParams(arg, param.Params))
		}

		// TODO: Extract parameters by parsing annotations defined in struct tags.
	}

	for name, p := range params {
		// Remain some unmatched entries.

		if !strings.HasPrefix(name, "__") {
			return fmt.Errorf("no argument `%s` declared in the method %s", name, method.Name)
		}

		// This is a blank identifier.
		typ := types.Typ[types.String] // Defaults to string.
		arg := &ifacetool.Param{
			Name:       name,
			TypeString: typ.Name(),
			Type:       typ,
		}
		op.Bind(arg, b.buildParams(arg, p.Params))
	}

	// Add path parameters according to the path pattern.
	for _, name := range extractPathVarNames(op.Pattern) {
		// If name is already bound to a path parameter by @kok(param) or
		// by struct tags, do not reset it.
		if isAlreadyPathParam(name, op.Request.Bindings) {
			continue
		}

		// Get the method argument name by always converting the path variable
		// name to lowerCamelCase (the naming convention in Go).
		//
		// Known issues:
		// - "xx_id" will be converted to "xxId" (not the conventional "xxID").
		argName := caseconv.ToLowerCamelCase(name)

		binding := op.Request.GetBinding(argName)
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
	if op.Request.BodyField == OptionNoBody {
		for _, binding := range op.Request.Bindings {
			if !binding.IsManual() {
				binding.SetIn(spec.InQuery)
			}
		}
	}

	return nil
}

func (b *OpBuilder) buildParams(arg *ifacetool.Param, annoParams []*spec.Parameter) []*spec.Parameter {
	name := caseconv.ToLowerCamelCase(arg.Name)
	if b.snakeCase {
		name = caseconv.ToSnakeCase(name)
	}

	defaultParam := &spec.Parameter{
		In:   spec.InBody, // Parameters are bound to the body by default.
		Name: name,
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
