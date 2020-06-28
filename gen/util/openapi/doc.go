package openapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/util/reflector"
)

var (
	reKok = regexp.MustCompile(`@kok\((\w+)\):\s*"(.+)"`)
)

func FromDoc(result *reflector.Result, doc map[string][]string) (*Specification, error) {
	spec := &Specification{}

	for _, m := range result.Interface.Methods {
		comments, ok := doc[m.Name]
		if !ok {
			return nil, fmt.Errorf("method %s has no comment", m.Name)
		}

		op := &Operation{Name: m.Name}

		// Add all request parameters with specified Name/Type
		params := make(map[string]*Param)
		for _, mp := range m.Params {
			p := &Param{
				In:   InBody, // param is in body by default
				Type: mp.Type,
			}
			p.SetName(mp.Name)
			op.addParam(p)

			// Build the mapping for later manipulation.
			params[p.Name] = p
		}

		// Set a default success response.
		op.Resp(http.StatusOK, MediaTypeJSON, nil)

		if err := manipulateByComments(op, params, comments); err != nil {
			return nil, err
		}

		spec.Operations = append(spec.Operations, op)
	}

	return spec, nil
}

func manipulateByComments(op *Operation, params map[string]*Param, comments []string) error {
	for _, comment := range comments {
		if !strings.Contains(comment, "@kok") {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 {
			return fmt.Errorf("invalid kok comment: %s", comment)
		}

		key, value := result[1], result[2]
		switch key {
		case "op":
			fields := strings.Fields(value)
			if len(fields) != 2 {
				return fmt.Errorf(`%q does not match the expected format: "<METHOD> <PATH>"`, value)
			}
			op.Method, op.Pattern = fields[0], fields[1]
		case "param":
			p := op.buildParam(value, "", "") // no default name and type

			name, subName := splitParamName(p.Name)
			param, ok := params[name]
			if !ok {
				return fmt.Errorf("no param `%s` declared in the method %s", name, op.Name)
			}

			if subName == "" {
				param.Set(p)
			} else {
				p.SetName(subName)
				param.Add(p)
			}
		case "success":
			op.SuccessResponse = buildResponse(value)
		case "errorEncoder":
			op.Options.ErrorEncoder = value
		default:
			return fmt.Errorf(`unrecognized kok key "%s" in comment: %s`, key, comment)
		}
	}

	if op.Method == "" {
		return fmt.Errorf("method %s has no comment about @kok(method)", op.Name)
	}

	if op.Pattern == "" {
		return fmt.Errorf("method %s has no comment about @kok(pattern)", op.Name)
	}

	return nil
}
