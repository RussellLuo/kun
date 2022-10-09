package parser

import (
	"fmt"
	"go/types"
	"regexp"

	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/parser"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

var (
	reCron = regexp.MustCompile(`^` + annotation.DirectiveCron.String() + `(.*)$`)
)

type Job struct {
	Name string
	Expr string
}

func Parse(data *ifacetool.Data, snakeCase bool) (map[string]*Job, error) {
	c := make(map[string]*Job)

	for _, m := range data.Methods {
		if !validateSignature(m) {
			return nil, fmt.Errorf("the signature of method %s must be `func(context.Context) error` when annotated by %s directive", m.Name, annotation.DirectiveCron)
		}

		for _, comment := range m.Doc {
			if annotation.Directive(comment).Dialect() != annotation.DialectCron {
				continue
			}

			result := reCron.FindStringSubmatch(comment)
			if len(result) != 2 {
				return nil, fmt.Errorf("invalid %s directive: %s", annotation.DirectiveCron, comment)
			}

			pairs, err := parser.ParseOptionPairs(result[1])
			if err != nil {
				return nil, err
			}

			job := new(Job)

			for _, pair := range pairs {
				switch pair.Key {
				case "name":
					job.Name = pair.Value
				case "expr":
					job.Expr = pair.Value
				default:
					return nil, fmt.Errorf(`unrecognized %s key "%s" in comment: %s`, annotation.DirectiveCron, pair.Key, comment)
				}
			}

			// Here we assume that all annotation keys are specified in the same line.

			if job.Name == "" {
				job.Name = caseconv.ToLowerCamelCase(m.Name)
				if snakeCase {
					job.Name = caseconv.ToSnakeCase(m.Name)
				}
			}

			if job.Expr == "" {
				return nil, fmt.Errorf(`missing key "expr" for %s directive in comment: %s`, annotation.DirectiveCron, comment)
			}

			c[m.Name] = job
			break // No need to continue since we have found the annotation.
		}
	}

	return c, nil
}

func validateSignature(m *ifacetool.Method) bool {
	if len(m.Params) != 1 || len(m.Returns) != 1 {
		return false
	}

	p := m.Params[0]
	if _, ok := p.Type.Underlying().(*types.Interface); !ok || p.TypeString != "context.Context" {
		return false
	}

	r := m.Returns[0]
	if _, ok := r.Type.Underlying().(*types.Interface); !ok || r.TypeString != "error" {
		return false
	}

	return true
}
