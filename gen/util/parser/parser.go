package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// Supported formats:
//   - k1=v1
//   - k2='v2'
var reOptionPair = regexp.MustCompile(`(\w+)=('[^']+'|\S+)`)

type OptionPair struct{ Key, Value string }

// ParseOptionPairs parses s into a list of OptionPair.
//
// Accepted formats:
//   - key1=value1
//   - key2='value2'
func ParseOptionPairs(s string) ([]OptionPair, error) {
	// NOTE: Instead of using ReplaceAllString and then FindAllStringSubmatch,
	// a more performant alternative solution may be to use ReplaceAllStringFunc once.

	unmatched := strings.TrimSpace(reOptionPair.ReplaceAllString(s, ""))
	if unmatched != "" {
		return nil, fmt.Errorf("invalid parameter option: %s", unmatched)
	}

	var pairs []OptionPair
	result := reOptionPair.FindAllStringSubmatch(s, -1)
	for _, r := range result {
		pairs = append(pairs, OptionPair{
			Key:   r[1],
			Value: strings.Trim(strings.TrimSpace(r[2]), "'"),
		})
	}
	return pairs, nil
}
