package docutil

import (
	"strings"
)

type Transport int

const (
	TransportHTTP Transport = 0b0001
	TransportGRPC Transport = 0b0010
	TransportAll  Transport = 0b0011
)

type Doc []string

func (d Doc) Transport() (t Transport) {
	for _, comment := range d {
		if IsKokGRPCAnnotation(comment) {
			t = t | TransportGRPC
		} else if IsKokAnnotation(comment) {
			t = t | TransportHTTP
		}
	}
	return t
}

// JoinComments joins backslash-continued comments.
func (d Doc) JoinComments() (joined Doc) {
	incompleteComment := ""

	for _, comment := range d {
		if incompleteComment == "" {
			if HasContinuationLine(comment) {
				incompleteComment = strings.TrimSuffix(comment, `\`)
			} else {
				joined = append(joined, comment)
			}
			continue
		}

		c := incompleteComment + strings.TrimSpace(comment)
		if HasContinuationLine(c) {
			incompleteComment = strings.TrimSuffix(c, `\`)
		} else {
			joined = append(joined, c)
			incompleteComment = ""
		}
	}

	return
}

func HasContinuationLine(comment string) bool {
	return strings.HasSuffix(comment, `\`)
}

func IsKokAnnotation(comment string) bool {
	return strings.HasPrefix(comment, "//kok:")
}

func IsKokGRPCAnnotation(comment string) bool {
	return strings.HasPrefix(comment, "//kok:grpc")
}
