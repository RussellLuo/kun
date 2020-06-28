package openapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func splitParamName(name string) (main string, sub string) {
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 1:
		// Non-nested parameter name.
		main, sub = parts[0], ""
	case 2:
		// Nested parameter name.
		main, sub = parts[0], parts[1]
	default:
		panic(fmt.Errorf("param name %q contains more than one `.`", name))
	}
	return
}

func isPrimitiveType(typ string) bool {
	switch typ {
	case "string", "bool",
		//"byte", "rune",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return true
	default:
		return false
	}
}

func buildResponse(text string) (resp *Response) {
	resp = new(Response)

	for _, part := range strings.Split(text, ",") {
		if !strings.Contains(part, ":") {
			panic(fmt.Errorf("invalid tag part: %s", part))
		}

		split := strings.SplitN(part, ":", 2)
		key, value := split[0], split[1]

		switch key {
		case "statusCode":
			var err error
			resp.StatusCode, err = strconv.Atoi(value)
			if err != nil {
				panic(fmt.Errorf("%q cannot be converted to an integer: %v", value, err))
			}
		case "encoder":
			resp.Options.Encoder = value
		default:
			panic(fmt.Errorf("invalid tag part: %s", part))
		}
	}

	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	if resp.MediaType == "" {
		resp.MediaType = MediaTypeJSON
	}

	return
}
