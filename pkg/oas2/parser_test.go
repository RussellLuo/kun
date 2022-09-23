package oas2

import (
	"reflect"
	"testing"
	"time"

	"github.com/RussellLuo/kun/pkg/httpcodec"
)

func TestParser_AddDefinition(t *testing.T) {
	toPtr := func(v interface{}) interface{} {
		switch t := v.(type) {
		// NOTE: t will be of type interface{}, if we list multiple types in each case.
		case bool:
			return &t
		case string:
			return &t
		case int:
			return &t
		case float32:
			return &t
		case float64:
			return &t
		}
		return &v
	}

	type Datum struct {
		Properties []string `json:"properties"`
	}

	type tree struct {
		value    string  // nolint:structcheck,unused
		parent   *tree   // nolint:structcheck,unused
		children []*tree // nolint:structcheck,unused
	}

	cases := []struct {
		name     string
		inBody   interface{}
		wantDefs map[string]Definition
	}{
		{
			name: "struct",
			inBody: struct {
				Name    string              `json:"name"`
				Male    bool                `json:"male"`
				Age     int                 `json:"age" kun:"descr=the-age"`
				Hobbies []string            `json:"hobbies"`
				Datum   *Datum              `json:"datum" kun:"descr=the-datum"`
				Data    []Datum             `json:"data"`
				Time    time.Time           `json:"time"`
				TimeStr string              `json:"time_str" kun:"type=time"`
				File    *httpcodec.FormFile `json:"file"`
				Other   map[string]bool     `json:"other" kun:"descr=other-is-a-map"`
			}{
				Name:    "xxx",
				Male:    true,
				Age:     10,
				Hobbies: []string{"music"},
				Datum:   nil,
				Data:    nil,
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
								Kind:        "basic",
								Type:        "integer",
								Format:      "int64",
								Description: "the-age",
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
							Name: "datum",
							Type: JSONType{
								Kind:        "object",
								Type:        "Datum",
								Description: "the-datum",
							},
						},
						{
							Name: "data",
							Type: JSONType{
								Kind: "array",
								Type: "Datum",
							},
						},
						{
							Name: "time",
							Type: JSONType{
								Kind:   "basic",
								Type:   "string",
								Format: "date-time",
							},
						},
						{
							Name: "time_str",
							Type: JSONType{
								Kind:   "basic",
								Type:   "string",
								Format: "date-time",
							},
						},
						{
							Name: "file",
							Type: JSONType{
								Kind:   "basic",
								Type:   "string",
								Format: "binary",
							},
						},
						{
							Name: "other",
							Type: JSONType{
								Kind:        "object",
								Type:        "Other",
								Description: "other-is-a-map",
							},
						},
					},
				},
				"Datum": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "properties",
							Type: JSONType{
								Kind: "array",
								Type: "string",
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
			name: "embedded struct",
			inBody: struct {
				Name string `json:"name"`
				Datum
			}{},
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
							Name: "properties",
							Type: JSONType{
								Kind: "array",
								Type: "string",
							},
						},
					},
				},
			},
		},
		{
			name: "embedded struct pointer",
			inBody: struct {
				Name string `json:"name"`
				*Datum
			}{},
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
							Name: "properties",
							Type: JSONType{
								Kind: "array",
								Type: "string",
							},
						},
					},
				},
			},
		},
		{
			name:   "array of string",
			inBody: []string{},
			wantDefs: map[string]Definition{
				"Response": {
					Type: "array",
					ItemTypeOrProperties: JSONType{
						Kind: "basic",
						Type: "string",
					},
				},
			},
		},
		{
			name:   "array of struct",
			inBody: []Datum{},
			wantDefs: map[string]Definition{
				"Response": {
					Type:                 "array",
					ItemTypeOrProperties: "Datum",
				},
				"Datum": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "properties",
							Type: JSONType{
								Kind: "array",
								Type: "string",
							},
						},
					},
				},
			},
		},
		{
			name: "map interface",
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
		{
			name:   "nil map interface",
			inBody: map[string]interface{}(nil),
			wantDefs: map[string]Definition{
				"Response": {
					Type:                 "object",
					ItemTypeOrProperties: []Property(nil),
				},
			},
		},
		{
			name:   "array of nil map interface",
			inBody: []map[string]interface{}(nil),
			wantDefs: map[string]Definition{
				"Response": {
					Type:                 "array",
					ItemTypeOrProperties: "ResponseArrayItem",
				},
				"ResponseArrayItem": {
					Type:                 "object",
					ItemTypeOrProperties: []Property(nil),
				},
			},
		},
		{
			name:   "array of array_of_nil_map_interface",
			inBody: [][]map[string]interface{}(nil),
			wantDefs: map[string]Definition{
				"Response": {
					Type:                 "array",
					ItemTypeOrProperties: "ResponseArrayItem",
				},
				"ResponseArrayItem": {
					Type:                 "array",
					ItemTypeOrProperties: "ResponseArrayItemArrayItem",
				},
				"ResponseArrayItemArrayItem": {
					Type:                 "object",
					ItemTypeOrProperties: []Property(nil),
				},
			},
		},
		{
			name: "pointers to basic types",
			inBody: struct {
				Name       *string  `json:"name"`
				Male       *bool    `json:"male"`
				Age        *int     `json:"age"`
				Float32Age *float32 `json:"float32_age"`
				Float64Age *float64 `json:"float64_age"`
			}{
				Name:       toPtr("xxx").(*string),
				Male:       toPtr(true).(*bool),
				Age:        toPtr(10).(*int),
				Float32Age: toPtr(float32(10)).(*float32),
				Float64Age: toPtr(10.0).(*float64),
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
							Name: "float32_age",
							Type: JSONType{
								Kind:   "basic",
								Type:   "number",
								Format: "float",
							},
						},
						{
							Name: "float64_age",
							Type: JSONType{
								Kind:   "basic",
								Type:   "number",
								Format: "double",
							},
						},
					},
				},
			},
		},
		{
			name:   "recursive type",
			inBody: new(tree),
			wantDefs: map[string]Definition{
				"Response": {
					Type: "object",
					ItemTypeOrProperties: []Property{
						{
							Name: "value",
							Type: JSONType{
								Kind: "basic",
								Type: "string",
							},
						},
						{
							Name: "parent",
							Type: JSONType{
								Kind: "object",
								Type: "Response",
							},
						},
						{
							Name: "children",
							Type: JSONType{
								Kind: "array",
								Type: "Response",
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewParser()
			p.AddDefinition("Response", reflect.ValueOf(c.inBody), false)
			defs := p.Definitions()
			if !reflect.DeepEqual(defs, c.wantDefs) {
				t.Fatalf("Defs: got (%#v), want (%#v)", defs, c.wantDefs)
			}
		})
	}
}
