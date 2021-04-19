package caseconv

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"already_snake", "already_snake"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}
	for _, c := range cases {
		got := ToSnakeCase(c.input)
		if got != c.want {
			t.Fatalf("Result: got (%#v), want (%#v)", got, c.want)
		}
	}
}

func TestToLowerCamelCase(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"alreadyCamel", "alreadyCamel"},
		{"a", "a"},
		{"aa_aa", "aaAa"},
		{"http_request", "httpRequest"},
		{"battery__life_Value", "batteryLifeValue"},
		{"id0_value", "id0Value"},
	}
	for _, c := range cases {
		got := ToLowerCamelCase(c.input)
		if got != c.want {
			t.Fatalf("Result: got (%#v), want (%#v)", got, c.want)
		}
	}
}
