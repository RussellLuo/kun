package oapi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RussellLuo/kok/reflector"
)

var (
	reKok = regexp.MustCompile(`@kok\((\w+)\):\s*"(.+)"`)
)

func FromDocs(docs []reflector.MethodDoc) (*Specification, error) {
	spec := &Specification{}

	for _, doc := range docs {
		op := &Operation{Name: doc.Name}

		for _, comment := range doc.Comments {
			if !strings.Contains(comment, "@kok") {
				continue
			}

			result := reKok.FindStringSubmatch(comment)
			if len(result) != 3 {
				return nil, fmt.Errorf("invalid kok comment: %s", comment)
			}

			key, value := result[1], result[2]
			switch key {
			case "path":
				op.Pattern = value
			case "method":
				op.Method = value
			case "param":
				op.addParam(op.buildParam(value, "", "")) // no default name and type
			default:
				return nil, fmt.Errorf("invalid kok comment: %s", comment)
			}
		}

		spec.Operations = append(spec.Operations, op)
	}

	return spec, nil
}
