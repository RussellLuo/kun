package annotation_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/gen/http/parser/annotation"
	"github.com/RussellLuo/kok/gen/http/spec"
)

func TestParseParam(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantOut    *annotation.Param
		wantErrStr string
	}{
		{
			name: "simple argument",
			in:   "name in=header name=X-User-Name required=true descr=user-name",
			wantOut: &annotation.Param{
				ArgName: "name",
				Params: []*spec.Parameter{
					{
						In:          spec.InHeader,
						Name:        "X-User-Name",
						Required:    true,
						Description: "user-name",
					},
				},
			},
		},
		{
			name: "argument aggregation",
			in:   "ip in=header name=X-Forwarded-For; in=request name=RemoteAddr",
			wantOut: &annotation.Param{
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
		{
			name: "no parameters",
			in:   "name",
			wantOut: &annotation.Param{
				ArgName: "name",
			},
		},
		{
			name: "in query by default",
			in:   "name required=true",
			wantOut: &annotation.Param{
				ArgName: "name",
				Params: []*spec.Parameter{
					{
						In:       spec.InQuery,
						Required: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, err := annotation.ParseParam(tt.in)
			if err != nil && err.Error() != tt.wantErrStr {
				t.Fatalf("ErrStr: got (%#v), want (%#v)", err.Error(), tt.wantErrStr)
			}
			if !reflect.DeepEqual(param, tt.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", param.Params[0], tt.wantOut.Params[0])
			}
		})
	}
}
