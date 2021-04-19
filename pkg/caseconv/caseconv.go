package caseconv

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

func ToCamelCase(s string) string {
	return matchUnderscore.ReplaceAllStringFunc(s, func(x string) string {
		return strings.Title(strings.TrimLeft(x, "_"))
	})
}

func ToUpperCamelCase(s string) string {
	return UpperFirst(ToCamelCase(s))
}

func ToLowerCamelCase(s string) string {
	return LowerFirst(ToCamelCase(s))
}

func UpperFirst(s string) string {
	return convertFirst(s, strings.ToUpper)
}

func LowerFirst(s string) string {
	return convertFirst(s, strings.ToLower)
}

func convertFirst(s string, f func(string) string) string {
	if len(s) == 0 {
		return s
	}
	return f(string(s[0])) + s[1:]
}
