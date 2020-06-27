package openapi

import (
	"fmt"
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
