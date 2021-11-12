package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/http/spec"
)

var (
	reKokParam = regexp.MustCompile(`^(\w+)(.*)$`)
)

type Param struct {
	ArgName string
	Params  []*spec.Parameter
}

// ParseParam parses s per the format as below:
//
//     <argName> [<parameter> [; <parameter2> [; ...]]]
//
// The format of `<parameter>`:
//
//     in=<in> name=<name> required=<required> type=<type> descr=<descr>
//
func ParseParam(s string) (*Param, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty //kok:param")
	}

	r := reKokParam.FindStringSubmatch(s)
	if len(r) != 3 {
		return nil, fmt.Errorf("invalid directive argument: %s", s)
	}
	argName, remaining := r[1], r[2]

	p := &Param{
		ArgName: argName,
	}

	if len(remaining) == 0 {
		// No remaining parameter definitions after the argument name.
		return p, nil
	}

	for _, text := range strings.Split(remaining, ";") {
		param, err := parseSingleParam(argName, text)
		if err != nil {
			return nil, err
		}
		p.Params = append(p.Params, param)
	}

	return p, nil
}

func parseSingleParam(argName, s string) (*spec.Parameter, error) {
	s = strings.TrimSpace(s)
	p := new(spec.Parameter)

	for _, part := range strings.Fields(s) {
		part = strings.TrimSpace(part)
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid directive argument: %s", part)
		}

		k, v := kv[0], kv[1]

		switch k {
		case "in":
			p.In = spec.Location(v)

			if err := errorIn(p.In); err != nil {
				return nil, err
			}
		case "name":
			p.Name = v
		case "required":
			p.Required = v == "true"
		case "type":
			p.Type = v
		case "descr":
			p.Description = v
		default:
			return nil, fmt.Errorf("invalid directive argument: %s", part)
		}
	}

	if p.In == spec.InPath {
		// This is a path parameter, thus it must be required.
		p.Required = true
	}

	if p.In == spec.InRequest && p.Name != "RemoteAddr" {
		return nil, fmt.Errorf("argument %q tries to extract value from `request.%s`, but only `request.RemoteAddr` is available", argName, p.Name)
	}

	if p.In == "" {
		// Defaults to be a query parameter.
		p.In = spec.InQuery
	}

	return p, nil
}

func errorIn(in spec.Location) error {
	if in != spec.InPath && in != spec.InQuery && in != spec.InHeader &&
		/*in != InCookie &&*/ in != spec.InRequest {

		return fmt.Errorf(
			"invalid location value: %s (must be %q, %q, %q or %q)",
			in, spec.InPath, spec.InQuery, spec.InHeader, spec.InRequest,
		)
	}
	return nil
}
