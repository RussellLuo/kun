package annotation_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/gen/http/parser/annotation"
)

func TestParseOp(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantOut    *annotation.Op
		wantErrStr string
	}{
		{
			name: "ok",
			in:   "POST /messages",
			wantOut: &annotation.Op{
				Method:  "POST",
				Pattern: "/messages",
			},
		},
		{
			name:       "invalid method",
			in:         "XXX /messages",
			wantErrStr: `invalid HTTP method "XXX" specified in "XXX /messages"`,
		},
		{
			name:       "invalid format",
			in:         "POST/messages",
			wantErrStr: `"POST/messages" does not match the expected format: <METHOD> <PATTERN>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := annotation.ParseOp(tt.in)
			if err != nil && err.Error() != tt.wantErrStr {
				t.Fatalf("ErrStr: got (%#v), want (%#v)", err.Error(), tt.wantErrStr)
			}
			if !reflect.DeepEqual(op, tt.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", op, tt.wantOut)
			}
		})
	}
}
