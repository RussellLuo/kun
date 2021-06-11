package openapi

import (
	"go/types"
	"reflect"
	"testing"
)

type structField struct {
	name string
	typ  types.Type
	tag  string
}

func newStruct(fields []*structField) *types.Struct {
	var fs []*types.Var
	var tags []string
	for _, f := range fields {
		fs = append(fs, types.NewField(0, nil, f.name, f.typ, false))
		tags = append(tags, f.tag)
	}
	return types.NewStruct(fs, tags)
}

func TestParser_Parse(t *testing.T) {
	cases := []struct {
		name            string
		inParser        *Parser
		inText          string
		wantAnnotations []*annotation
		wantErrStr      string
	}{
		{
			name: "int argument",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name:    "arg1",
						RawType: types.Typ[types.Int],
					},
				},
			},
			inText: "arg1 < in:query",
			wantAnnotations: []*annotation{
				{
					ArgName: "arg1",
					In:      "query",
					Type:    "int",
				},
			},
		},
		{
			name: "int argument with name",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name:    "arg1",
						RawType: types.Typ[types.Int],
					},
				},
			},
			inText: "arg1 < in:query,name:arg",
			wantAnnotations: []*annotation{
				{
					ArgName: "arg1",
					In:      "query",
					Name:    "arg",
					Type:    "int",
				},
			},
		},
		{
			name: "bool argument",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name:    "arg1",
						RawType: types.Typ[types.Bool],
					},
				},
			},
			inText: "arg1 < in:header",
			wantAnnotations: []*annotation{
				{
					ArgName: "arg1",
					In:      "header",
					Type:    "bool",
				},
			},
		},
		{
			name: "required bool argument",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name:    "arg1",
						RawType: types.Typ[types.Bool],
					},
				},
			},
			inText: "arg1 < in:header,required:true",
			wantAnnotations: []*annotation{
				{
					ArgName:  "arg1",
					In:       "header",
					Type:     "bool",
					Required: true,
				},
			},
		},
		{
			name: "struct argument",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name: "arg1",
						RawType: newStruct([]*structField{
							{
								name: "Field1",
								typ:  types.Typ[types.Int],
							},
							{
								name: "Field2",
								typ:  types.Typ[types.Uint],
								tag:  `kok:".field2"`,
							},
							{
								name: "Field3",
								typ:  types.Typ[types.Bool],
								tag:  `kok:"header.field3"`,
							},
							{
								name: "Field4",
								typ:  types.Typ[types.String],
								tag:  `kok:"path.field4,required"`,
							},
							{
								name: "Field5",
								typ:  types.Typ[types.String],
								tag:  `kok:"-"`,
							},
						}),
					},
				},
			},
			inText: "arg1",
			wantAnnotations: []*annotation{
				{
					ArgName: "arg1",
					In:      "query",
					Name:    "Field1",
					Type:    "int",
				},
				{
					ArgName: "arg1",
					In:      "query",
					Name:    "field2",
					Type:    "uint",
				},
				{
					ArgName: "arg1",
					In:      "header",
					Name:    "field3",
					Type:    "bool",
				},
				{
					ArgName:  "arg1",
					In:       "path",
					Name:     "field4",
					Type:     "string",
					Required: true,
				},
			},
		},
		{
			name: "invalid in",
			inParser: &Parser{
				methodName: "Test",
				params: map[string]*Param{
					"arg1": {
						Name:    "arg1",
						RawType: types.Typ[types.Int],
					},
				},
			},
			inText:     "arg1 < in:xxx",
			wantErrStr: `invalid location value: xxx (must be "path", "query", "header" or "request")`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			annotations, err := c.inParser.Parse(c.inText)
			if (err == nil && c.wantErrStr != "") || (err != nil && err.Error() != c.wantErrStr) {
				t.Fatalf("Err: got (%#v), want (%#v)", err, c.wantErrStr)
			}
			if !reflect.DeepEqual(annotations, c.wantAnnotations) {
				t.Fatalf("Annotations: got (%#v), want (%#v)", annotations, c.wantAnnotations)
			}
		})
	}
}
