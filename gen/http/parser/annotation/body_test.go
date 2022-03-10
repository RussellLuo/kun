package annotation_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kun/gen/http/parser/annotation"
)

func TestParseBody(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantOut    *annotation.Body
		wantErrStr string
	}{
		{
			name: "field",
			in:   "user",
			wantOut: &annotation.Body{
				Field: "user",
			},
		},
		{
			name: "manipulation",
			in:   "user name=u type=string descr=user-name",
			wantOut: &annotation.Body{
				Manipulations: map[string]*annotation.Manipulation{
					"user": {
						Name:        "u",
						Type:        "string",
						Description: "user-name",
					},
				},
			},
		},
		{
			name:       "in unsupported",
			in:         "user in=path",
			wantErrStr: "parameter `in` is unsupported in body manipulation",
		},
		{
			name:       "required unsupported",
			in:         "user required=true",
			wantErrStr: "parameter `required` is unsupported in body manipulation",
		},
		{
			name:       "invalid directive",
			in:         "user name:xx",
			wantErrStr: "invalid parameter option: name:xx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := annotation.ParseBody(tt.in)
			if err != nil && err.Error() != tt.wantErrStr {
				t.Fatalf("ErrStr: got (%#v), want (%#v)", err.Error(), tt.wantErrStr)
			}
			if !reflect.DeepEqual(body, tt.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", body, tt.wantOut)
			}
		})
	}
}
