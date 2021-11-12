package annotation

import (
	"fmt"
	"strings"

	"github.com/RussellLuo/kok/gen/http/spec"
)

// ParseMetadata parses doc per the format as below:
//
//     <property>=<value>
//
func ParseMetadata(doc []string) (*spec.Metadata, error) {
	m := &spec.Metadata{
		DocsPath:    "/api",
		Title:       "No Title",
		Version:     "0.0.0",
		Description: GetDescriptionFromDoc(doc),
		BasePath:    "/",
	}

	for _, comment := range doc {
		if !IsKokAnnotation(comment) {
			continue
		}

		result := reKok.FindStringSubmatch(comment)
		if len(result) != 3 || result[1] != "oas" {
			if result[1] == "alias" {
				continue
			}
			return nil, fmt.Errorf("invalid kok directive: %s", comment)
		}

		value := strings.TrimSpace(result[2])
		kv := strings.SplitN(value, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf(`%q does not match the expected format: "<key>=<value>"`, value)
		}

		k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		switch k {
		case "docsPath":
			m.DocsPath = v
		case "title":
			m.Title = v
		case "version":
			m.Version = v
		case "basePath":
			m.BasePath = v
		case "tags":
			m.DefaultTags = strings.Split(v, ",")
		default:
			return nil, fmt.Errorf(`invalid key %q for //kok:oas in %q`, k, value)
		}
	}

	return m, nil
}

func GetDescriptionFromDoc(doc []string) string {
	var comments []string
	for _, comment := range doc {
		if !IsKokAnnotation(comment) {
			comments = append(comments, strings.TrimPrefix(comment, "// "))
		}
	}
	// Separate multiline description by raw `\n`.
	return strings.Join(comments, "\\n")
}
