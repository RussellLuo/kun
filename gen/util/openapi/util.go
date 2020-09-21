package openapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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

func buildSuccessResponse(text string) (resp *Response, encoder string) {
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
			encoder = value
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

func getFailureResponseEncoder(text string) string {
	for _, part := range strings.Split(text, ",") {
		if !strings.Contains(part, ":") {
			panic(fmt.Errorf("invalid tag part: %s", part))
		}

		split := strings.SplitN(part, ":", 2)
		key, value := split[0], split[1]

		switch key {
		case "encoder":
			return value
		default:
			panic(fmt.Errorf("invalid tag part: %s", part))
		}
	}
	panic(fmt.Errorf("empty @kok(failure)"))
}
