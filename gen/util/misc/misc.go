package misc

import (
	"regexp"
	"strings"
)

var (
	// matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(s string) string {
	// snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake := matchAllCap.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(string(s[0])) + s[1:]
}
