package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kun/gen/http/spec"
	"github.com/RussellLuo/kun/gen/util/annotation"
)

const (
	OptionNoBody = "-"
)

var (
	reSingleVarName = regexp.MustCompile(`^\w+$`)
)

type Manipulation struct {
	Name        string
	Type        string
	Description string
}

type Body struct {
	Field         string
	Manipulations map[string]*Manipulation
}

// ParseBody parses s per the format as below:
//
//     <field> or <manipulation> [; <manipulation2> [; ...]]
//
// The format of `<manipulation>`:
//
//     <argName> name=<name> type=<type> descr=<descr>
//
func ParseBody(s string) (*Body, error) {
	// Simple format: <field>
	if s == OptionNoBody || reSingleVarName.MatchString(s) {
		return &Body{Field: s}, nil
	}

	// Complicated format: <manipulation> [; <manipulation2> [; ...]]
	m := make(map[string]*Manipulation)
	for _, text := range strings.Split(s, ";") {
		text = strings.TrimSpace(text)
		if text == "" {
			break
		}

		param, err := parseParam(text)
		if err != nil {
			return nil, err
		}

		if len(param.Params) != 1 {
			return nil, fmt.Errorf("bad manipulation %q in %s", s, annotation.DirectiveHTTPBody)
		}
		p := param.Params[0]

		if p.In != spec.InQuery {
			// XXX: Handle the case of manual definition `in=query`.
			return nil, fmt.Errorf("parameter `in` is unsupported in body manipulation")
		}
		if p.Required {
			return nil, fmt.Errorf("parameter `required` is unsupported in body manipulation")
		}

		m[param.ArgName] = &Manipulation{
			Name:        p.Name,
			Type:        p.Type,
			Description: p.Description,
		}
	}

	if len(m) > 0 {
		return &Body{Manipulations: m}, nil
	}

	return nil, fmt.Errorf("invalid %s directive: %s", annotation.DirectiveHTTPBody, s)
}
