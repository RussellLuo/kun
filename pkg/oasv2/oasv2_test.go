package oasv2

import (
	"reflect"
	"testing"
)

func TestAddDefinition(t *testing.T) {
	cases := []struct {
		name     string
		inBody   interface{}
		wantDefs map[string]Definition
	}{
		{
			name: "struct",
			inBody: struct {
				Name    string          `json:"name"`
				Male    bool            `json:"male"`
				Age     int             `json:"age"`
				Hobbies []string        `json:"hobbies"`
				Other   map[string]bool `json:"other"`
			}{
				Name:    "xxx",
				Male:    true,
				Age:     10,
				Hobbies: []string{"music"},
				Other:   map[string]bool{"married": true},
			},
			wantDefs: map[string]Definition{
				"Response": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "name",
							Type: JSONType{
								Kind: "basic",
								Type: "string",
							},
						},
						{
							Name: "male",
							Type: JSONType{
								Kind: "basic",
								Type: "boolean",
							},
						},
						{
							Name: "age",
							Type: JSONType{
								Kind:   "basic",
								Type:   "integer",
								Format: "int64",
							},
						},
						{
							Name: "hobbies",
							Type: JSONType{
								Kind: "array",
								Type: "string",
							},
						},
						{
							Name: "other",
							Type: JSONType{
								Kind: "object",
								Type: "Other",
							},
						},
					},
				},
				"Other": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "married",
							Type: JSONType{
								Kind: "basic",
								Type: "boolean",
							},
						},
					},
				},
			},
		},
		{
			name: "map_interface",
			inBody: map[string]interface{}{
				"attrs": map[string]interface{}{
					"age": 20,
				},
			},
			wantDefs: map[string]Definition{
				"Response": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "attrs",
							Type: JSONType{
								Kind: "object",
								Type: "Attrs",
							},
						},
					},
				},
				"Attrs": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "age",
							Type: JSONType{
								Kind:   "basic",
								Type:   "integer",
								Format: "int64",
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defs := make(map[string]Definition)
			AddDefinition(defs, "Response", reflect.ValueOf(c.inBody))
			if !reflect.DeepEqual(defs, c.wantDefs) {
				t.Fatalf("Defs: got (%#v), want (%#v)", defs, c.wantDefs)
			}
		})
	}
}
