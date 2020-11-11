package openapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/RussellLuo/kok/gen/util/reflector"
)

func buildSuccessResponse(text string, results map[string]*reflector.Param, opName string) *Response {
	resp := new(Response)

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
		case "body":
			if _, ok := results[value]; !ok {
				panic(fmt.Errorf("no result `%s` declared in the method %s", value, opName))
			}
			resp.BodyField = value
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

	return resp
}
