package httpclient

import (
	"go/types"
	"testing"

	"github.com/RussellLuo/kun/pkg/ifacetool"
)

func TestEmptyValue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input *ifacetool.Param
		want  string
	}{
		{
			name: "int",
			input: &ifacetool.Param{
				Name:       "param",
				TypeString: "int",
				Type:       types.Typ[types.Int],
			},
			want: "0",
		},
		{
			name: "string",
			input: &ifacetool.Param{
				Name:       "param",
				TypeString: "string",
				Type:       types.Typ[types.String],
			},
			want: `""`,
		},
		{
			name: "bool",
			input: &ifacetool.Param{
				Name:       "param",
				TypeString: "bool",
				Type:       types.Typ[types.Bool],
			},
			want: "false",
		},
		{
			name: "struct",
			input: &ifacetool.Param{
				Name:       "param",
				TypeString: "MyStruct",
				Type:       &types.Struct{},
			},
			want: "MyStruct{}",
		},
		{
			name: "iface",
			input: &ifacetool.Param{
				Name:       "param",
				TypeString: "MyInterface",
				Type:       &types.Interface{},
			},
			want: "nil",
		},
	} {
		t.Run("test empty value: " + tt.name, func(t *testing.T) {
			got := emptyValue(tt.input)
			if got != tt.want {
				t.Fatalf("got (%v), want (%#v)", got, tt.want)
			}
		})
	}
}
