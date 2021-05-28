package openapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/util/reflector"
	"github.com/RussellLuo/kok/pkg/caseconv"
)

const (
	OptionNoBody = "-"
)

var (
	reKok = regexp.MustCompile(`@kok\((\w+)\):\s*(.+)$`)

	rePathVarName = regexp.MustCompile(`{(\w+)}`)

	reSingleVarName = regexp.MustCompile(`^\w+$`)
)

func FromDoc(result *reflector.Result, doc *reflector.InterfaceDoc) (*Specification, error) {
	metadata, err := buildMetadata(doc.Doc)
	if err != nil {
		return nil, err
	}

	spec := &Specification{
		Metadata: metadata,
	}

	for _, m := range result.Interface.Methods {
		comments, ok := doc.MethodDocs[m.Name]
		if !ok || !hasKokAnnotations(comments) {
			continue
		}

		op := &Operation{
			Name:        m.Name,
			Description: getDescriptionFromDoc(comments),
		}

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

func getDescriptionFromDoc(doc []string) string {
	var comments []string
	for _, comment := range doc {
		if !isKokAnnotation(comment) {
			comments = append(comments, strings.TrimPrefix(comment, "// "))
		}
	}
	// Separate multiline description by raw `\n`.
	return strings.Join(comments, "\\n")
}

func hasKokAnnotations(comments []string) bool {
	for _, comment := range comments {
		if isKokAnnotation(comment) {
			return true
		}
	}
	return false
}

func isKokAnnotation(comment string) bool {
	content := strings.TrimPrefix(comment, "//")
	trimmed := strings.TrimSpace(content)
	return strings.HasPrefix(trimmed, "@kok")
}

func manipulateByComments(op *Operation, params map[string]*Param, results map[string]*reflector.Param, comments []string) error {
	parser := &Parser{
		methodName: op.Name,
		params:     params,
	}

	getParam := func(argName string) (*Param, error) {
		p, ok := params[argName]
		if !ok {
			return nil, fmt.Errorf("no param `%s` declared in the method %s", argName, op.Name)
		}
		return p, nil
	}

	setParamByAnnotation := func(a *annotation) error {
		param, err := getParam(a.ArgName)
		if err != nil {
			return err
		}

		if !param.inUse {
			param.SetByAnnotation(a)
		} else {
			copied := *param
			param.SetByAnnotation(a)

			// Add a new parameter with the same name.
			op.addParam(&copied)
		}

		return nil
	}

	for _, comment := range comments {
		if !isKokAnnotation(comment) {
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
			annotations, err := parser.Parse(value)
			if err != nil {
				return err
			}
			for _, a := range annotations {
				if err := setParamByAnnotation(a); err != nil {
					return err
				}
			}

		case "body":
			setBodyField := func(value string) error {
				if _, ok := params[value]; value != OptionNoBody && !ok {
					return fmt.Errorf("no param `%s` declared in the method %s", value, op.Name)
				}
				op.Request.BodyField = value
				return nil
			}

			// Simple format: <field>
			if reSingleVarName.MatchString(value) {
				if err := setBodyField(value); err != nil {
					return err
				}
				continue
			}

			// Complicated format:
			//
			// body:<field>,name:<argName>=<NAME>,descr:<argName>=<DESCR>

			getParamAndValue := func(v string) (*Param, string, error) {
				s := strings.SplitN(v, "=", 2)
				if len(s) != 2 {
					return nil, "", fmt.Errorf(`%q does not match the expected format: "<argName>=<value>"`, v)
				}
				argName, value := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])

				param, err := getParam(argName)
				if err != nil {
					return nil, "", err
				}

				return param, value, nil
			}

			for _, subValue := range strings.Split(value, ",") {
				parts := strings.SplitN(subValue, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf(`%q does not match the expected format: "<key>:<value>"`, subValue)
				}
				k, v := parts[0], parts[1]

				switch k {
				case "body":
					if err := setBodyField(v); err != nil {
						return err
					}
				case "name":
					param, alias, err := getParamAndValue(v)
					if err != nil {
						return err
					}
					param.SetAlias(alias)
				case "descr":
					param, descr, err := getParamAndValue(v)
					if err != nil {
						return err
					}
					param.SetDescription(descr)
				default:
					return fmt.Errorf("invalid tag part: %s", subValue)
				}
			}

		case "success":
			op.SuccessResponse = buildSuccessResponse(value, results, op.Name)

		case "oas":
			parts := strings.SplitN(value, ":", 2)
			if len(parts) != 2 || parts[0] != "tags" {
				return fmt.Errorf(`%q does not match the expected format: "tags:<tag1>[,<tag2>]"`, value)
			}
			op.Tags = strings.Split(parts[1], ",")

		default:
			return fmt.Errorf(`unrecognized kok key "%s" in comment: %s`, key, comment)
		}
	}

	if op.Method == "" && op.Pattern == "" {
		return fmt.Errorf("method %s has no comment about @kok(op)", op.Name)
	}

	// Add path parameters according to the path pattern.
	for _, name := range extractPathVarNames(op.Pattern) {
		// If name is already bound to a path parameter by @kok(param) or
		// by struct tags, do not reset it.
		if isAlreadyPathParam(name, op.Request.Params) {
			continue
		}

		// Add this path parameter.
		annotations, err := parser.Parse(name + " < in:path")
		if err != nil {
			return err
		}
		for _, a := range annotations {
			if err := setParamByAnnotation(a); err != nil {
				return err
			}
		}
	}

	// Add possible query parameters if no-request-body is specified.
	if op.Request.BodyField == OptionNoBody {
		annotations, err := makeKokQueryParamTextsFromBodyParams(parser)
		if err != nil {
			return err
		}
		for _, a := range annotations {
			if err := setParamByAnnotation(a); err != nil {
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
		name := caseconv.ToLowerCamelCase(s[1])
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

func makeKokQueryParamTextsFromBodyParams(parser *Parser) (annotations []*annotation, err error) {
	for _, p := range parser.params {
		if p.In != InBody || p.Name == "ctx" {
			// Ignore non-body parameters and the special context.Context parameter.
			continue
		}

		annos, err := parser.Parse(p.Name + " < in:query")
		if err != nil {
			return nil, err
		}

		annotations = append(annotations, annos...)
	}
	return
}
