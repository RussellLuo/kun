package openapi

import (
	"fmt"
	"go/types"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/util/misc"
	"github.com/RussellLuo/kok/gen/util/reflector"
	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"
)

const (
	OptionNoBody = "-"
)

var (
	reKok = regexp.MustCompile(`@kok2?\((\w+)\):\s*(.+)$`)

	rePathVarName = regexp.MustCompile(`{(\w+)}`)
)

func FromDoc(result *reflector.Result, doc map[string][]string) (*Specification, error) {
	spec := &Specification{}

	for _, m := range result.Interface.Methods {
		comments, ok := doc[m.Name]
		if !ok {
			continue
		}

		op := &Operation{Name: m.Name}

		// Add all request parameters with specified Name/Type
		params := make(map[string]*Param)
		for _, mp := range m.Params {
			p := &Param{
				In:        InBody, // param is in body by default
				Type:      mp.Type,
				RawType:   mp.RawType, // used for adding query parameters later
				AliasType: mp.Type,
			}
			p.SetName(mp.Name)
			op.addParam(p)

			// Build the mapping for later manipulation.
			params[p.Name] = p
		}

		results := make(map[string]*reflector.Param)
		for _, mr := range m.Returns {
			results[mr.Name] = mr
		}

		// Set a default success response.
		op.Resp(http.StatusOK, MediaTypeJSON, nil)

		if err := manipulateByComments(op, params, results, comments); err != nil {
			return nil, err
		}

		spec.Operations = append(spec.Operations, op)
	}

	return spec, nil
}

func manipulateByComments(op *Operation, params map[string]*Param, results map[string]*reflector.Param, comments []string) error {
	var prevParamName string

	setParam := func(value, prevName string) (string, error) {
		p := op.buildParamV2(value, prevName)

		param, ok := params[p.Name]
		if !ok {
			return "", fmt.Errorf("no param `%s` declared in the method %s", p.Name, op.Name)
		}

		if !param.inUse {
			param.Set(p)
		} else {
			copied := *param
			param.Set(p)

			// Add a new parameter with the same name.
			op.addParam(&copied)
		}

		return p.Name, nil
	}

	for _, comment := range comments {
		if !strings.Contains(comment, "@kok") {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 {
			return fmt.Errorf("invalid kok comment: %s", comment)
		}

		key, value := result[1], strings.TrimSpace(result[2])
		switch key {
		case "op":
			fields := strings.Fields(value)
			if len(fields) != 2 {
				return fmt.Errorf(`%q does not match the expected format: <METHOD> <PATTERN>`, value)
			}
			op.Method, op.Pattern = fields[0], fields[1]

		case "param":
			name, err := setParam(value, prevParamName)
			if err != nil {
				return err
			}
			prevParamName = name

		case "body":
			if _, ok := params[value]; value != OptionNoBody && !ok {
				return fmt.Errorf("no param `%s` declared in the method %s", value, op.Name)
			}
			op.Request.BodyField = value

		case "success":
			op.SuccessResponse = buildSuccessResponse(value, results, op.Name)

		default:
			return fmt.Errorf(`unrecognized kok key "%s" in comment: %s`, key, comment)
		}
	}

	if op.Method == "" && op.Pattern == "" {
		return fmt.Errorf("method %s has no comment about @kok(op)", op.Name)
	}

	// Add path parameters according to the path pattern.
	for _, name := range extractPathVarNames(op.Pattern) {
		// If name is already bound to a path parameter that is specified in
		// @kok(param), do not reset it.
		if isAlreadyPathParam(name, op.Request.Params) {
			continue
		}

		// Build the @kok(param) value according to the path variable name.
		text := fmt.Sprintf("%s < in:path", name)

		// Add this path parameter.
		if _, err := setParam(text, ""); err != nil {
			return err
		}
	}

	// Add possible query parameters if no-request-body is specified.
	if op.Request.BodyField == OptionNoBody {
		for _, text := range makeKokQueryParamTextsFromBodyParams(params) {
			if _, err := setParam(text, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func extractPathVarNames(pattern string) (names []string) {
	result := rePathVarName.FindAllStringSubmatch(pattern, -1)
	if len(result) == 0 {
		return
	}

	for _, s := range result {
		// Convert possible snake case to camel case.
		name := misc.ToLowerCamelCase(s[1])
		names = append(names, name)
	}
	return
}

func isAlreadyPathParam(name string, params []*Param) bool {
	for _, param := range params {
		if param.In == InPath && param.Alias == name {
			return true
		}
	}
	return false
}

func makeKokQueryParamTextsFromBodyParams(params map[string]*Param) (texts []string) {
	for _, p := range params {
		if p.In != InBody || p.Name == "ctx" {
			// Ignore non-body parameters and the special context.Context parameter.
			continue
		}

		nameTypes := make(map[string]string)

		// Build the kok tag text.
		switch t := p.RawType.Underlying().(type) {
		case *types.Basic:
			nameTypes[p.Alias] = t.Name()
		case *types.Slice:
			et, ok := t.Elem().(*types.Basic)
			if !ok {
				panic(fmt.Errorf("query parameter cannot be mapped to argument %q of type %v", p.Name, t))
			}
			nameTypes[p.Alias] = "[]" + et.Name()
		case *types.Struct:
			for i := 0; i < t.NumFields(); i++ {
				var typeName string
				switch ft := t.Field(i).Type().(type) {
				case *types.Basic:
					typeName = ft.Name()
				case *types.Slice:
					et, ok := ft.Elem().(*types.Basic)
					if !ok {
						panic(fmt.Errorf("query parameter cannot be mapped to struct field %q of type %v", t.Field(i).Name(), ft))
					}
					typeName = "[]" + et.Name()
				default:
					panic(fmt.Errorf("query parameter cannot be mapped to struct field %q of type %v", t.Field(i).Name(), ft))
				}

				field := reflect.StructField{
					Tag:  reflect.StructTag(t.Tag(i)),
					Name: t.Field(i).Name(),
				}
				name, _, _ := httpcodec.GetFieldName(field)
				if strings.HasPrefix(name, "query.") {
					// Only add the field that is mapped to a query parameter.
					name = strings.TrimPrefix(name, "query.") // get the real query name
					nameTypes[name] = typeName
				}
			}
		default:
			panic(fmt.Errorf("query parameter cannot be mapped to argument %q of type %v", p.Name, t))
		}

		for name, typ := range nameTypes {
			texts = append(texts, fmt.Sprintf("%s < in:query,name:%s,type:%s", p.Name, name, typ))
		}
	}

	return
}
