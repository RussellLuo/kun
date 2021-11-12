package annotation

import (
	"fmt"
	"regexp"
)

const (
	OptionNoBody = "-"
)

var (
	reSingleVarName = regexp.MustCompile(`^\w+$`)
)

type Manipulation struct {
	Name        string
	Type        string
	Description string
}

type Body struct {
	Field         string
	Manipulations []*Manipulation
}

// ParseBody parses s per the format as below:
//
//     <field> or body=<field> manip=`<manipulation> [; <manipulation2> [; ...]]`
//
// The format of `<manipulation>`:
//
//     <argName> name=<name> type=<type> descr=<descr>
//
func ParseBody(s string) (*Body, error) {
	// Simple format: <field>
	if s == OptionNoBody || reSingleVarName.MatchString(s) {
		return &Body{Field: s}, nil
	}

	// Complicated format: <manipulation> [; <manipulation2> [; ...]]
	// TODO: add support for parsing complicated format.

	return nil, fmt.Errorf("invalid //kok:body directive: %s", s)
}
