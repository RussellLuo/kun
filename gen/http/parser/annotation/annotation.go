package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kun/gen/http/spec"
	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

var (
	reHTTP = regexp.MustCompile(`^` + annotation.DirectivePrefix + `(\w+)\s+(.+)$`)
)

type InterfaceAnnotation struct {
	Metadata *spec.Metadata
	Aliases  Aliases
}

func ParseInterfaceAnnotation(doc []string) (*InterfaceAnnotation, error) {
	m, err := ParseMetadata(doc)
	if err != nil {
		return nil, err
	}

	aliases, err := ParseAliases(doc)
	if err != nil {
		return nil, err
	}

	return &InterfaceAnnotation{
		Metadata: m,
		Aliases:  aliases,
	}, nil
}

type MethodAnnotation struct {
	Op      *Op
	Params  map[string]*Param
	Body    *Body
	Success *spec.Response
	Tags    []string
}

func ParseMethodAnnotation(method *ifacetool.Method) (*MethodAnnotation, error) {
	anno := &MethodAnnotation{Params: make(map[string]*Param)}

	for _, comment := range method.Doc {
		if annotation.Directive(comment).Dialect() != annotation.DialectHTTP {
			continue
		}

		result := reHTTP.FindStringSubmatch(comment)
		if len(result) != 3 {
			return nil, fmt.Errorf("invalid %s directive: %s", annotation.Name, comment)
		}

		key, value := result[1], strings.TrimSpace(result[2])
		switch d := annotation.FromSubDirective(key); d {
		case annotation.DirectiveHTTPOp:
			if anno.Op != nil {
				return nil, fmt.Errorf("duplicate %s directive in: %s", d, comment)
			}

			op, err := ParseOp(value)
			if err != nil {
				return nil, err
			}
			anno.Op = op

		case annotation.DirectiveHTTPParam:
			params, err := ParseParams(value)
			if err != nil {
				return nil, err
			}
			for _, p := range params {
				anno.Params[p.ArgName] = p
			}

		case annotation.DirectiveHTTPBody:
			if anno.Body != nil {
				return nil, fmt.Errorf("duplicate %s directive in: %s", d, comment)
			}

			body, err := ParseBody(value)
			if err != nil {
				return nil, err
			}
			anno.Body = body

		case annotation.DirectiveHTTPSuccess:
			if anno.Success != nil {
				return nil, fmt.Errorf("duplicate %s directive in: %s", d, comment)
			}

			success, err := ParseSuccess(value, method)
			if err != nil {
				return nil, err
			}
			anno.Success = success

		case annotation.DirectiveHTTPOAS:
			if len(anno.Tags) > 0 {
				return nil, fmt.Errorf("duplicate %s directive in: %s", d, comment)
			}

			parts := strings.SplitN(value, ":", 2)
			if len(parts) != 2 || parts[0] != "tags" {
				return nil, fmt.Errorf(`%q does not match the expected format: "tags:<tag1>[,<tag2>]"`, value)
			}
			anno.Tags = strings.Split(parts[1], ",")

		default:
			return nil, fmt.Errorf(`unrecognized %s directive "%s" in comment: %s`, annotation.Name, key, comment)
		}
	}

	return anno, nil
}
