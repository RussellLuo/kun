package annotation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/gen/http/spec"
	"github.com/RussellLuo/kok/gen/util/docutil"
	"github.com/RussellLuo/kok/pkg/ifacetool"
)

var (
	reKok = regexp.MustCompile(`^//kok:(\w+)\s+(.+)$`)
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
		if !docutil.IsKokAnnotation(comment) || docutil.IsKokGRPCAnnotation(comment) {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 {
			return nil, fmt.Errorf("invalid kok directive: %s", comment)
		}

		key, value := result[1], strings.TrimSpace(result[2])
		switch key {
		case "op":
			if anno.Op != nil {
				return nil, fmt.Errorf("duplicate //kok:op directive in: %s", comment)
			}

			op, err := ParseOp(value)
			if err != nil {
				return nil, err
			}
			anno.Op = op

		case "param":
			params, err := ParseParams(value)
			if err != nil {
				return nil, err
			}
			for _, p := range params {
				anno.Params[p.ArgName] = p
			}

		case "body":
			if anno.Body != nil {
				return nil, fmt.Errorf("duplicate //kok:body directive in: %s", comment)
			}

			body, err := ParseBody(value)
			if err != nil {
				return nil, err
			}
			anno.Body = body

		case "success":
			if anno.Success != nil {
				return nil, fmt.Errorf("duplicate //kok:success directive in: %s", comment)
			}

			success, err := ParseSuccess(value, method)
			if err != nil {
				return nil, err
			}
			anno.Success = success

		case "oas":
			if len(anno.Tags) > 0 {
				return nil, fmt.Errorf("duplicate //kok:oas directive in: %s", comment)
			}

			parts := strings.SplitN(value, ":", 2)
			if len(parts) != 2 || parts[0] != "tags" {
				return nil, fmt.Errorf(`%q does not match the expected format: "tags:<tag1>[,<tag2>]"`, value)
			}
			anno.Tags = strings.Split(parts[1], ",")

		default:
			return nil, fmt.Errorf(`unrecognized kok key "%s" in comment: %s`, key, comment)
		}
	}

	return anno, nil
}
