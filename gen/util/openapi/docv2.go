package openapi

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reKokV2 = regexp.MustCompile(`@kok2\((\w+)\):\s*(.+)$`)
)

func manipulateByCommentsV2(op *Operation, params map[string]*Param, comments []string) error {
	prevParamName := ""

	for _, comment := range comments {
		if !strings.Contains(comment, "@kok2(") {
			continue
		}

		result := reKokV2.FindStringSubmatch(comment)
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
			p := op.buildParamV2(value, prevParamName)
			prevParamName = p.Name

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
			op.SuccessResponse, _ = buildSuccessResponse(value)
		default:
			return fmt.Errorf(`unrecognized kok key "%s" in comment: %s`, key, comment)
		}
	}

	return nil
}
