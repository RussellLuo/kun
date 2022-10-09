package parser

import (
	"fmt"
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
		job := new(Job)
		c[m.Name] = job

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

			if job.Name == "" {
				job.Name = caseconv.ToLowerCamelCase(m.Name)
				if snakeCase {
					job.Name = caseconv.ToSnakeCase(m.Name)
				}
			}

			if job.Expr == "" {
				return nil, fmt.Errorf(`missing key "expr" for %s directive in comment: %s`, annotation.DirectiveCron, comment)
			}
		}
	}

	return c, nil
}
