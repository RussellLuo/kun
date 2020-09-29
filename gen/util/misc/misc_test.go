package misc

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
