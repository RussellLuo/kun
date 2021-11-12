package annotation

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/RussellLuo/kok/gen/http/spec"
	"github.com/RussellLuo/kok/pkg/ifacetool"
)

// ParseSuccess parses s per the format as below:
//
//     statusCode=<statusCode> body=<body> manip=`<manipulation> [; <manipulation2> [; ...]]`
//
// The format of `<manipulation>`:
//
//     <argName> name=<name> type=<type> descr=<descr>
//
func ParseSuccess(s string, method *ifacetool.Method) (*spec.Response, error) {
	resp := new(spec.Response)

	returns := make(map[string]*ifacetool.Param)
	for _, r := range method.Returns {
		returns[r.Name] = r
	}

	for _, part := range strings.Fields(s) {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid argument format: %s", part)
		}

		key, value := kv[0], kv[1]

		switch key {
		case "statusCode":
			var err error
			resp.StatusCode, err = strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("%q cannot be converted to an integer: %v", value, err)
			}
		case "body":
			if _, ok := returns[value]; !ok {
				return nil, fmt.Errorf("no result `%s` declared in the method %s", value, method.Name)
			}
			resp.BodyField = value
		case "manip":
			// TODO: Add support for the `manip` argument.
			fallthrough
		default:
			return nil, fmt.Errorf("invalid tag part: %s", part)
		}
	}

	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	if resp.MediaType == "" {
		resp.MediaType = spec.MediaTypeJSON
	}

	return resp, nil
}
