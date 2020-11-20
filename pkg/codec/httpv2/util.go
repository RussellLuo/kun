package codec

import (
	"strings"
)

const (
	sep = ","
)

func QueryListToString(l []string) string {
	var b strings.Builder
	for _, v := range l {
		if v != "" {
			b.WriteString(v + sep)
		}
	}
	s := b.String()
	return strings.TrimRight(s, sep)
}

func QueryStringToList(s string) (l []string) {
	for _, v := range strings.Split(s, sep) {
		if v != "" {
			l = append(l, v)
		}
	}
	return
}
