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

func buildMetadata(comments []string) (m *Metadata, err error) {
	m = &Metadata{
		Title:       "No Title",
		Version:     "0.0.0",
		Description: getDescriptionFromDoc(comments),
		BasePath:    "/",
	}

	for _, comment := range comments {
		if !isKokAnnotation(comment) {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 || result[1] != "oas" {
			return nil, fmt.Errorf("invalid kok comment: %s", comment)
		}

		value := strings.TrimSpace(result[2])
		parts := strings.SplitN(value, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(`%q does not match the expected format: "<key>:<value>"`, value)
		}

		k, v := parts[0], parts[1]
		switch k {
		case "title":
			m.Title = v
		case "version":
			m.Version = v
		case "basePath":
			m.BasePath = v
		case "tags":
			m.DefaultTags = strings.Split(v, ",")
		default:
			return nil, fmt.Errorf(`invalid key %q for @kok(oas) in %q`, k, value)
		}
	}

	return m, nil
}
