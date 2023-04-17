package parser

import (
	"encoding/json"
	"go/types"
	"reflect"
	"testing"

	"github.com/RussellLuo/kun/gen/http/parser/annotation"
	"github.com/RussellLuo/kun/gen/http/spec"
	"github.com/RussellLuo/kun/pkg/ifacetool"
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

func TestOpBuilder_Build(t *testing.T) {
	tests := []struct {
		name          string
		inOpBuilder   *OpBuilder
		inMethod      *ifacetool.Method
		wantOperation []*spec.Operation
		wantErrStr    string
	}{
		{
			name:        `multiple routes for one method`,
			inOpBuilder: &OpBuilder{snakeCase: true},
			inMethod: &ifacetool.Method{
				Name: "MultipleRoutesForOneMethod",
				Params: []*ifacetool.Param{
					{
						Name:       "token",
						TypeString: "string",
						Type:       types.Typ[types.String],
						Variadic:   false,
					},
				},
				Doc: []string{
					"//kun:op GET /route_one",
					"//kun:op POST /route_two",
					"//kun:param token in=header name=Authorization required=true",
				},
			},
			wantOperation: []*spec.Operation{
				{
					Name:         `MultipleRoutesForOneMethod`,
					GoMethodName: `MultipleRoutesForOneMethod`,
					Method:       `GET`,
					Pattern:      `/route_one`,
					Request: &spec.Request{
						Bindings: []*spec.Binding{
							{
								Arg: &ifacetool.Param{
									Name:       "token",
									TypeString: "string",
									Variadic:   false,
								},
								Params: []*spec.Parameter{
									{
										In:       "header",
										Name:     "Authorization",
										Required: true,
										Type:     "string",
									},
								},
							},
						},
					},
					SuccessResponse: &spec.Response{
						StatusCode: 200,
						MediaType:  `application/json; charset=utf-8`,
					},
				},
				{
					Name:         `MultipleRoutesForOneMethod1`,
					GoMethodName: `MultipleRoutesForOneMethod`,
					Method:       `POST`,
					Pattern:      `/route_two`,
					Request: &spec.Request{
						Bindings: []*spec.Binding{
							{
								Arg: &ifacetool.Param{
									Name:       "token",
									TypeString: "string",
									Variadic:   false,
								},
								Params: []*spec.Parameter{
									{
										In:       "header",
										Name:     "Authorization",
										Required: true,
										Type:     "string",
									},
								},
							},
						},
					},
					SuccessResponse: &spec.Response{
						StatusCode: 200,
						MediaType:  `application/json; charset=utf-8`,
					},
				},
			},
			wantErrStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.inOpBuilder.Build(tt.inMethod)
			if (err == nil && tt.wantErrStr != "") || (err != nil && err.Error() != tt.wantErrStr) {
				t.Fatalf("Err: got (%#v), want (%#v)", err, tt.wantErrStr)
			}
			if len(res) != len(tt.wantOperation) {
				t.Fatalf("Len unmatch: got (%#v), want (%#v)", len(res), len(tt.wantOperation))
			}
			for i, operation := range res {
				d1, _ := json.Marshal(operation)
				d2, _ := json.Marshal(tt.wantOperation[i])

				if !reflect.DeepEqual(d1, d2) {
					t.Fatalf("Operation(%d): got (%#v), want (%#v)", i, res[i], tt.wantOperation[i])
				}
			}
		})
	}
}

func TestOpBuilder_setParams(t *testing.T) {
	tests := []struct {
		name           string
		inOpBuilder    *OpBuilder
		inReq          *spec.Request
		inMethod       *ifacetool.Method
		inParams       map[string]*annotation.Param
		inPathVarNames []string
		wantBindings   map[string][]*spec.Parameter
		wantErrStr     string
	}{
		{
			name:        "normal bindings",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "arg1",
						TypeString: "int",
						Type:       types.Typ[types.Int],
					},
					{
						Name:       "arg2",
						TypeString: "string",
						Type:       types.Typ[types.String],
					},
					{
						Name:       "arg3",
						TypeString: "StructArg",
						Type: newStruct([]*structField{
							{
								name: "Field1",
								typ:  types.Typ[types.Uint],
							},
							{
								name: "Field2",
								typ:  types.Typ[types.Bool],
								tag:  `kun:"in=header name=field_2"`,
							},
						}),
					},
					{
						Name:       "arg4",
						TypeString: "int",
						Type:       types.Typ[types.Int],
					},
				},
			},
			inParams: map[string]*annotation.Param{
				"arg1": {
					ArgName: "arg1",
					Params: []*spec.Parameter{
						{
							In: spec.InPath,
						},
					},
				},
				"arg2": {
					ArgName: "arg2",
					Params: []*spec.Parameter{
						{
							In:   spec.InHeader,
							Name: "X-Name",
						},
					},
				},
				"arg3": {
					ArgName: "arg3",
				},
			},
			wantBindings: map[string][]*spec.Parameter{
				"arg1": {
					{
						In:   spec.InPath,
						Name: "arg1",
						Type: "int",
					},
				},
				"arg2": {
					{
						In:   spec.InHeader,
						Name: "X-Name",
						Type: "string",
					},
				},
				"arg3": {
					{
						In:   spec.InQuery,
						Name: "field1",
						Type: "uint",
					},
					{
						In:   spec.InHeader,
						Name: "field_2",
						Type: "bool",
					},
				},
				"arg4": {
					{
						In:   spec.InBody,
						Name: "arg4",
						Type: "int",
					},
				},
			},
		},
		{
			name:        "blank identifier",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name:   "Test",
				Params: nil,
			},
			inParams: map[string]*annotation.Param{
				"__": {
					ArgName: "__",
					Params: []*spec.Parameter{
						{
							In:       spec.InHeader,
							Name:     "Authorization",
							Required: true,
						},
					},
				},
			},
			wantBindings: map[string][]*spec.Parameter{
				"__": {
					{
						In:       spec.InHeader,
						Name:     "Authorization",
						Type:     "string",
						Required: true,
					},
				},
			},
		},
		{
			name:        "path auto-binding",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "id",
						TypeString: "int",
						Type:       types.Typ[types.Int],
					},
				},
			},
			inParams:       nil,
			inPathVarNames: []string{"id"},
			wantBindings: map[string][]*spec.Parameter{
				"id": {
					{
						In:       spec.InPath,
						Name:     "id",
						Type:     "int",
						Required: true,
					},
				},
			},
		},
		{
			name:        "nobody to query (basic arguments)",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       &spec.Request{BodyField: "-"},
			inMethod: &ifacetool.Method{
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "arg1",
						TypeString: "int",
						Type:       types.Typ[types.Int],
					},
					{
						Name:       "arg2",
						TypeString: "string",
						Type:       types.Typ[types.String],
					},
				},
			},
			inParams: nil,
			wantBindings: map[string][]*spec.Parameter{
				"arg1": {
					{
						In:   spec.InQuery,
						Name: "arg1",
						Type: "int",
					},
				},
				"arg2": {
					{
						In:   spec.InQuery,
						Name: "arg2",
						Type: "string",
					},
				},
			},
		},
		{
			name:        "nobody to query (struct argument)",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       &spec.Request{BodyField: "-"},
			inMethod: &ifacetool.Method{
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "arg1",
						TypeString: "StructArg",
						Type: newStruct([]*structField{
							{
								name: "Field1",
								typ:  types.Typ[types.Uint],
							},
							{
								name: "Field2",
								typ:  types.Typ[types.Bool],
							},
							{
								name: "Field3",
								typ:  types.Typ[types.String],
								tag:  `kun:"in=header"`,
							},
						}),
					},
				},
			},
			inParams: nil,
			wantBindings: map[string][]*spec.Parameter{
				"arg1": {
					{
						In:   spec.InQuery,
						Name: "field1",
						Type: "uint",
					},
					{
						In:   spec.InQuery,
						Name: "field2",
						Type: "bool",
					},
					{
						In:   spec.InHeader,
						Name: "field3",
						Type: "string",
					},
				},
			},
		},
		{
			name:        "non-basic argument",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "arg1",
						TypeString: "int",
						Type:       types.NewSlice(types.NewSlice(types.Typ[types.Int])),
					},
				},
			},
			inParams: map[string]*annotation.Param{
				"arg1": {
					ArgName: "arg1",
					Params: []*spec.Parameter{
						{
							In: spec.InQuery,
						},
					},
				},
			},
			wantBindings: map[string][]*spec.Parameter{},
			wantErrStr:   `cannot define extra parameters for non-basic argument "arg1"`,
		},
		{
			name:        "unmatched annotation parameter",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name:   "Test",
				Params: nil,
			},
			inParams: map[string]*annotation.Param{
				"arg1": {
					ArgName: "arg1",
					Params: []*spec.Parameter{
						{
							In: spec.InQuery,
						},
					},
				},
			},
			wantBindings: map[string][]*spec.Parameter{},
			wantErrStr:   `no argument "arg1" declared in the method Test`,
		},
		{
			name:        "unsuccessful path auto-binding",
			inOpBuilder: &OpBuilder{snakeCase: true},
			inReq:       new(spec.Request),
			inMethod: &ifacetool.Method{
				Name:   "Test",
				Params: nil,
			},
			inParams:       nil,
			inPathVarNames: []string{"id"},
			wantBindings:   map[string][]*spec.Parameter{},
			wantErrStr:     `cannot bind path parameter "id": no argument "id" declared in the method Test`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inOpBuilder.setParams(tt.inReq, tt.inMethod, tt.inParams, tt.inPathVarNames)
			if (err == nil && tt.wantErrStr != "") || (err != nil && err.Error() != tt.wantErrStr) {
				t.Fatalf("Err: got (%#v), want (%#v)", err, tt.wantErrStr)
			}

			gotBindings := make(map[string][]*spec.Parameter)
			for _, b := range tt.inReq.Bindings {
				gotBindings[b.Arg.Name] = b.Params
			}
			if !reflect.DeepEqual(gotBindings, tt.wantBindings) {
				t.Fatalf("Bindings: got (%#v), want (%#v)", gotBindings, tt.wantBindings)
			}
		})
	}
}

func TestOpBuilder_inferAnnotationParams(t *testing.T) {
	tests := []struct {
		name         string
		inOpBuilder  *OpBuilder
		inMethodName string
		inArg        *ifacetool.Param
		wantParams   []*spec.Parameter
		wantErrStr   string
	}{
		{
			name:         "basic argument",
			inOpBuilder:  &OpBuilder{snakeCase: true},
			inMethodName: "Test",
			inArg: &ifacetool.Param{
				Name:       "arg1",
				TypeString: "int",
				Type:       types.Typ[types.Int],
			},
			wantParams: []*spec.Parameter{
				{
					In:   spec.InQuery,
					Name: "arg1",
					Type: "int",
				},
			},
		},
		{
			name:         "slice argument",
			inOpBuilder:  &OpBuilder{snakeCase: true},
			inMethodName: "Test",
			inArg: &ifacetool.Param{
				Name:       "arg1",
				TypeString: "[]bool",
				Type:       types.NewSlice(types.Typ[types.Bool]),
			},
			wantParams: []*spec.Parameter{
				{
					In:   spec.InQuery,
					Name: "arg1",
					Type: "[]bool",
				},
			},
		},
		{
			name:         "struct argument",
			inOpBuilder:  &OpBuilder{snakeCase: true},
			inMethodName: "Test",
			inArg: &ifacetool.Param{
				Name:       "arg1",
				TypeString: "int",
				Type: newStruct([]*structField{
					{
						name: "Field1",
						typ:  types.Typ[types.Int],
					},
					{
						name: "Field2",
						typ:  types.Typ[types.Uint],
						tag:  `kun:"name=field_2"`,
					},
					{
						name: "Field3",
						typ:  types.Typ[types.Bool],
						tag:  `kun:"in=header name=field_3"`,
					},
					{
						name: "Field4",
						typ:  types.Typ[types.String],
						tag:  `kun:"in=path name=field_4 required=true"`,
					},
					{
						name: "Field5",
						typ:  types.Typ[types.String],
						tag:  `kun:"name=-"`,
					},
				}),
			},
			wantParams: []*spec.Parameter{
				{
					In:   spec.InQuery,
					Name: "field1",
					Type: "int",
				},
				{
					In:   spec.InQuery,
					Name: "field_2",
					Type: "uint",
				},
				{
					In:   spec.InHeader,
					Name: "field_3",
					Type: "bool",
				},
				{
					In:       spec.InPath,
					Name:     "field_4",
					Type:     "string",
					Required: true,
				},
			},
		},
		{
			name:         "basic argument",
			inOpBuilder:  &OpBuilder{snakeCase: true},
			inMethodName: "Test",
			inArg: &ifacetool.Param{
				Name:       "arg1",
				TypeString: "int",
				Type:       types.NewSlice(types.NewSlice(types.Typ[types.String])),
			},
			wantErrStr: `parameter cannot be mapped to argument "arg1" (of type [][]string) in method Test`,
		},
		{
			name:         "interface argument",
			inOpBuilder:  &OpBuilder{snakeCase: true},
			inMethodName: "Test",
			inArg: &ifacetool.Param{
				Name:       "ctx",
				TypeString: "context.Context",
				Type:       types.NewInterfaceType(nil, nil),
			},
			wantErrStr: "", // no error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := tt.inOpBuilder.inferAnnotationParams(tt.inMethodName, tt.inArg)
			if (err == nil && tt.wantErrStr != "") || (err != nil && err.Error() != tt.wantErrStr) {
				t.Fatalf("Err: got (%#v), want (%#v)", err, tt.wantErrStr)
			}
			if !reflect.DeepEqual(params, tt.wantParams) {
				t.Fatalf("Params: got (%#v), want (%#v)", params, tt.wantParams)
			}
		})
	}
}

func Test_extractPathVarNames(t *testing.T) {
	want := []string{"name", "path"}
	got := extractPathVarNames("/var/{name}/in/{path}")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Names: got (%#v), want (%#v)", got, want)
	}
}

func Test_removePathParamsNotItsOwn(t *testing.T) {
	tests := []struct {
		name           string
		inBindings     []*spec.Binding
		inPathVarNames []string
		wantBindings   []*spec.Binding
	}{
		{
			name: "normal binding",
			inBindings: []*spec.Binding{
				{
					Arg: &ifacetool.Param{
						Name:       "userID",
						TypeString: "string",
						Type:       types.Typ[types.String],
					},
					Params: []*spec.Parameter{
						{
							In:       spec.InPath,
							Name:     "userID",
							Required: true,
							Type:     "string",
						},
					},
				},
				{
					Arg: &ifacetool.Param{
						Name:       "messageID",
						TypeString: "string",
						Type:       types.Typ[types.String],
					},
					Params: []*spec.Parameter{
						{
							In:       spec.InPath,
							Name:     "messageID",
							Required: true,
							Type:     "string",
						},
					},
				},
			},
			inPathVarNames: []string{"messageID"},
			wantBindings: []*spec.Binding{
				{
					Arg: &ifacetool.Param{
						Name:       "messageID",
						TypeString: "string",
						Type:       types.Typ[types.String],
					},
					Params: []*spec.Parameter{
						{
							In:       spec.InPath,
							Name:     "messageID",
							Required: true,
							Type:     "string",
						},
					},
				},
			},
		},
		{
			name: "aggregate binding",
			inBindings: []*spec.Binding{
				{
					Arg: &ifacetool.Param{
						Name:       "req",
						TypeString: "GetMessageRequest",
						Type:       nil,
					},
					Params: []*spec.Parameter{
						{
							In:       spec.InPath,
							Name:     "userID",
							Required: true,
							Type:     "string",
						},
						{
							In:       spec.InPath,
							Name:     "messageID",
							Required: true,
							Type:     "string",
						},
					},
				},
			},
			inPathVarNames: []string{"messageID"},
			wantBindings: []*spec.Binding{
				{
					Arg: &ifacetool.Param{
						Name:       "req",
						TypeString: "GetMessageRequest",
						Type:       nil,
					},
					Params: []*spec.Parameter{
						{
							In:       spec.InPath,
							Name:     "messageID",
							Required: true,
							Type:     "string",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBindings := removePathParamsNotItsOwn(tt.inBindings, tt.inPathVarNames)
			if !reflect.DeepEqual(gotBindings, tt.wantBindings) {
				t.Fatalf("Bindings: got (%#v), want (%#v)", gotBindings, tt.wantBindings)
			}
		})
	}
}

func TestStructField_Parse(t *testing.T) {
	tests := []struct {
		name        string
		in          *StructField
		wantOmitted bool
		wantParams  []*spec.Parameter
		wantErrStr  string
	}{
		{
			name: "in query by default",
			in: &StructField{
				Name: "Name",
				Type: "string",
			},
			wantOmitted: false,
			wantParams: []*spec.Parameter{
				{
					In:   spec.InQuery,
					Name: "name",
					Type: "string",
				},
			},
		},
		{
			name: "in path",
			in: &StructField{
				Name: "Name",
				Type: "string",
				Tag:  `kun:"in=path"`,
			},
			wantOmitted: false,
			wantParams: []*spec.Parameter{
				{
					In:       spec.InPath,
					Name:     "name",
					Required: true,
					Type:     "string",
				},
			},
		},
		{
			name: "omitted",
			in: &StructField{
				Name: "Name",
				Type: "string",
				Tag:  `kun:"name=-"`,
			},
			wantOmitted: true,
			wantParams:  nil,
		},
		{
			name: "required",
			in: &StructField{
				Name: "Name",
				Type: "string",
				Tag:  `kun:"required=true"`,
			},
			wantOmitted: false,
			wantParams: []*spec.Parameter{
				{
					In:       spec.InQuery,
					Name:     "name",
					Required: true,
					Type:     "string",
				},
			},
		},
		{
			name: "has type",
			in: &StructField{
				Name: "Name",
				Type: "string",
				Tag:  `kun:"type=bool"`,
			},
			wantOmitted: false,
			wantParams: []*spec.Parameter{
				{
					In:   spec.InQuery,
					Name: "name",
					Type: "bool",
				},
			},
		},
		{
			name: "has description",
			in: &StructField{
				Name: "Name",
				Type: "string",
				Tag:  `kun:"descr=the-description"`,
			},
			wantOmitted: false,
			wantParams: []*spec.Parameter{
				{
					In:          spec.InQuery,
					Name:        "name",
					Type:        "string",
					Description: "the-description",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Parse()
			if (err == nil && tt.wantErrStr != "") || (err != nil && err.Error() != tt.wantErrStr) {
				t.Fatalf("Err: got (%#v), want (%#v)", err, tt.wantErrStr)
			}

			if tt.in.Omitted != tt.wantOmitted {
				t.Fatalf("Omitted: got (%#v), want (%#v)", tt.in.Omitted, tt.wantOmitted)
			}

			if !reflect.DeepEqual(tt.in.Params, tt.wantParams) {
				t.Fatalf("Params: got (%#v), want (%#v)", tt.in.Params, tt.wantParams)
			}
		})
	}
}
