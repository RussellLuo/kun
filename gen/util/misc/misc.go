package misc

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

	matchUnderscore = regexp.MustCompile("(_+[0-9A-Za-z])")
)

func ToSnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToLowerCamelCase(s string) string {
	camel := matchUnderscore.ReplaceAllStringFunc(s, func(x string) string {
		return strings.Title(strings.TrimLeft(x, "_"))
	})
	return LowerFirst(camel)
}

func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
