package annotation_test

import (
	"github.com/RussellLuo/kun/gen/http/spec"
	"go/types"
	"reflect"
	"testing"

	"github.com/RussellLuo/kun/gen/http/parser/annotation"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

func TestParseMethodAnnotation(t *testing.T) {
	tests := []struct {
		name       string
		inMethod   *ifacetool.Method
		inAliases  annotation.Aliases
		wantOut    *annotation.MethodAnnotation
		wantErrStr string
	}{
		{
			name: "using aliases in //kun:param",
			inMethod: &ifacetool.Method{
				Doc:  []string{"//kun:param $opID"},
				Name: "Test",
				Params: []*ifacetool.Param{
					{
						Name:       "operatorID",
						TypeString: "int",
						Type:       types.Typ[types.Int],
					},
				},
			},
			inAliases: annotation.Aliases{
				"opID": `operatorID in=header name=Authorization required=true`,
			},
			wantOut: &annotation.MethodAnnotation{
				Params: map[string]*annotation.Param{
					"operatorID": {
						ArgName: "operatorID",
						Params: []*spec.Parameter{
							{
								In:       spec.InHeader,
								Name:     "Authorization",
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anno, err := annotation.ParseMethodAnnotation(tt.inMethod, tt.inAliases)
			if err != nil && err.Error() != tt.wantErrStr {
				t.Fatalf("ErrStr: got (%#v), want (%#v)", err.Error(), tt.wantErrStr)
			}
			if !reflect.DeepEqual(anno, tt.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", anno, tt.wantOut)
			}
		})
	}
}
