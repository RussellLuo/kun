package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

var (
	reEvent = regexp.MustCompile(`^` + annotation.DirectiveEvent.String() + `(.*)$`)
)

type EventInfo struct {
	Types     map[string]string // method names => event types
	DataField string
}

func Parse(data *ifacetool.Data, snakeCase bool) (*EventInfo, error) {
	e := &EventInfo{Types: make(map[string]string)}

	for _, m := range data.Methods {
		for _, comment := range m.Doc {
			if annotation.Directive(comment).Dialect() != annotation.DialectEvent {
				continue
			}

			result := reEvent.FindStringSubmatch(comment)
			if len(result) != 2 {
				return nil, fmt.Errorf("invalid %s directive: %s", annotation.DirectiveEvent, comment)
			}

			value := strings.TrimSpace(result[1])
			if value == "" {
				e.Types[m.Name] = caseconv.ToLowerCamelCase(m.Name)
				if snakeCase {
					e.Types[m.Name] = caseconv.ToSnakeCase(m.Name)
				}
				continue
			}

			fields := strings.Fields(value)
			for _, f := range fields {
				parts := strings.Split(f, "=")
				if len(parts) != 2 {
					return nil, fmt.Errorf(`%q does not match the expected format: <key>=<value>`, f)
				}
				k, v := parts[0], parts[1]

				switch k {
				case "type":
					e.Types[m.Name] = v
				case "data":
					if !isMethodParam(m, v) {
						return nil, fmt.Errorf("no argument %q declared in the method %s", v, m.Name)
					}
					e.DataField = v
				default:
					return nil, fmt.Errorf(`unrecognized %s key "%s" in comment: %s`, annotation.Name, k, comment)
				}
			}
		}
	}

	return e, nil
}

func isMethodParam(m *ifacetool.Method, name string) bool {
	for _, p := range m.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}
