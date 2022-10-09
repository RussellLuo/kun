package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kun/gen/http/spec"
	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/parser"
)

var (
	reKunParam = regexp.MustCompile(`^(\w+)(.*)$`)
)

type Param struct {
	ArgName string
	Params  []*spec.Parameter
}

// ParseParams parses s per the format as below:
//
//	<argName> [<parameter> [, <parameter2> [, ...]]]
//
// The format of `<parameter>`:
//
//	in=<in> name=<name> required=<required> type=<type> descr=<descr>
//
// Multiple bindings can be specified in a single semicolon-separated comment.
func ParseParams(s string) ([]*Param, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty %s", annotation.DirectiveHTTPParam)
	}

	var list []*Param

	for _, text := range strings.Split(s, ";") {
		p, err := parseParam(text)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

func parseParam(s string) (*Param, error) {
	s = strings.TrimSpace(s)

	r := reKunParam.FindStringSubmatch(s)
	if len(r) != 3 {
		return nil, fmt.Errorf("invalid directive arguments: %s", s)
	}
	argName, remaining := r[1], strings.TrimSpace(r[2])

	p := &Param{ArgName: argName}

	if remaining == "" {
		// No remaining parameter definitions after the argument name.
		return p, nil
	}

	opts, err := ParseParamOptions(argName, remaining)
	if err != nil {
		return nil, err
	}
	p.Params = append(p.Params, opts...)

	return p, nil
}

func ParseParamOptions(argName, s string) ([]*spec.Parameter, error) {
	var params []*spec.Parameter
	for _, text := range strings.Split(s, ",") {
		param, err := parseOption(argName, strings.TrimSpace(text))
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

func parseOption(argName, s string) (*spec.Parameter, error) {
	s = strings.TrimSpace(s)
	p := new(spec.Parameter)

	pairs, err := parser.ParseOptionPairs(s)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		switch pair.Key {
		case "in":
			p.In = spec.Location(pair.Value)

			if err := validateLocation(p.In); err != nil {
				return nil, err
			}
		case "name":
			p.Name = pair.Value
		case "required":
			p.Required = pair.Value == "true"
		case "type":
			p.Type = pair.Value
		case "descr":
			p.Description = pair.Value
		default:
			return nil, fmt.Errorf("invalid parameter option: %s=%s", pair.Key, pair.Value)
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
		// Location defaults to query if not specified.
		p.In = spec.InQuery
	}

	return p, nil
}

func validateLocation(in spec.Location) error {
	if in != spec.InPath && in != spec.InQuery && in != spec.InHeader &&
		/*in != InCookie &&*/ in != spec.InRequest {

		return fmt.Errorf(
			"invalid location value: %s (must be %q, %q, %q or %q)",
			in, spec.InPath, spec.InQuery, spec.InHeader, spec.InRequest,
		)
	}
	return nil
}
