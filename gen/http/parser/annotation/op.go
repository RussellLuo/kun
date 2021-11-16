package annotation

import (
	"fmt"
	"net/http"
	"strings"
)

type Op struct {
	Method  string
	Pattern string
}

// ParseOp parses s per the format as below:
//
//     <method> <pattern>
//
func ParseOp(s string) (*Op, error) {
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return nil, fmt.Errorf("%q does not match the expected format: <METHOD> <PATTERN>", s)
	}
	method, pattern := fields[0], fields[1]

	if _, ok := httpMethods[method]; !ok {
		return nil, fmt.Errorf("invalid HTTP method %q specified in %q", method, s)
	}

	return &Op{
		Method:  method,
		Pattern: pattern,
	}, nil
}

var (
	httpMethods = map[string]struct{}{
		http.MethodOptions: {},
		http.MethodHead:    {},
		http.MethodGet:     {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
	}
)
