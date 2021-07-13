package openapi

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	reKokAlias = regexp.MustCompile("(\\w+)\\s*=\\s*`([^`]+)`$")
	reAliasVar = regexp.MustCompile(`\$(\w+)`)
)

type Aliases map[string]string

func newAliases(comments []string) (Aliases, error) {
	a := make(map[string]string)

	for _, comment := range comments {
		if !isKokAnnotation(comment) {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 || result[1] != "alias" {
			if result[1] == "oas" {
				continue
			}
			return nil, fmt.Errorf("invalid kok comment: %s", comment)
		}

		value := strings.TrimSpace(result[2])
		r := reKokAlias.FindStringSubmatch(value)
		if len(r) != 3 {
			return nil, fmt.Errorf("%q does not match the expected format: \"<name>=`<value>`\"", value)
		}
		k, v := r[1], r[2]

		a[k] = v
	}

	return Aliases(a), nil
}

func (a Aliases) Set(key, value string) {
	a[key] = value
}

// Eval replaces all possible aliases with their values.
func (a Aliases) Eval(value string) (string, error) {
	var err error
	return reAliasVar.ReplaceAllStringFunc(value, func(s string) string {
		k := strings.TrimPrefix(s, "$")
		v, ok := a[k]
		if !ok {
			err = fmt.Errorf("unknown alias %q", s)
		}
		return v
	}), err
}
