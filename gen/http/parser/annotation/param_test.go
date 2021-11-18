package annotation_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/gen/http/parser/annotation"
	"github.com/RussellLuo/kok/gen/http/spec"
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantOut    []*annotation.Param
		wantErrStr string
	}{
		{
			name: "one binding one sub-parameter",
			in:   "name in=header name=X-User-Name required=true type=string descr=user-name",
			wantOut: []*annotation.Param{
				{
					ArgName: "name",
					Params: []*spec.Parameter{
						{
							In:          spec.InHeader,
							Name:        "X-User-Name",
							Required:    true,
							Type:        "string",
							Description: "user-name",
						},
					},
				},
			},
		},
		{
			name: "one binding no sub-parameter",
			in:   "name",
			wantOut: []*annotation.Param{
				{
					ArgName: "name",
				},
			},
		},
		{
			name: "one binding multiple sub-parameters",
			in:   "ip in=header name=X-Forwarded-For, in=request name=RemoteAddr",
			wantOut: []*annotation.Param{
				{
					ArgName: "ip",
					Params: []*spec.Parameter{
						{
							In:   spec.InHeader,
							Name: "X-Forwarded-For",
						},
						{
							In:   spec.InRequest,
							Name: "RemoteAddr",
						},
					},
				},
			},
		},
		{
			name: "multiple bindings",
			in:   "name; age in=query; ip in=header name=X-Forwarded-For, in=request name=RemoteAddr",
			wantOut: []*annotation.Param{
				{
					ArgName: "name",
				},
				{
					ArgName: "age",
					Params: []*spec.Parameter{
						{
							In: spec.InQuery,
						},
					},
				},
				{
					ArgName: "ip",
					Params: []*spec.Parameter{
						{
							In:   spec.InHeader,
							Name: "X-Forwarded-For",
						},
						{
							In:   spec.InRequest,
							Name: "RemoteAddr",
						},
					},
				},
			},
		},
		{
			name: "defaults in query",
			in:   "name required=true",
			wantOut: []*annotation.Param{
				{
					ArgName: "name",
					Params: []*spec.Parameter{
						{
							In:       spec.InQuery,
							Required: true,
						},
					},
				},
			},
		},
		{
			name:       "invalid parameter pair",
			in:         "name x:y",
			wantErrStr: "invalid parameter pair: x:y",
		},
		{
			name:       "invalid location",
			in:         "name in=xxx",
			wantErrStr: `invalid location value: xxx (must be "path", "query", "header" or "request")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := annotation.ParseParams(tt.in)
			if err != nil && err.Error() != tt.wantErrStr {
				t.Fatalf("ErrStr: got (%#v), want (%#v)", err.Error(), tt.wantErrStr)
			}
			if !reflect.DeepEqual(list, tt.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", list, tt.wantOut)
			}
		})
	}
}
