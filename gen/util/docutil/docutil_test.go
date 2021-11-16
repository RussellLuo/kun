package docutil_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/kok/gen/util/docutil"
)

func TestDoc_JoinLines(t *testing.T) {
	tests := []struct {
		name string
		in   docutil.Doc
		want docutil.Doc
	}{
		{
			name: "no backslash",
			in: []string{
				"//kok:op POST /logs",
				"//kok:param ip in=header name=X-Forwarded-For; in=request name=RemoteAddr",
			},
			want: []string{
				"//kok:op POST /logs",
				"//kok:param ip in=header name=X-Forwarded-For; in=request name=RemoteAddr",
			},
		},
		{
			name: "has backslash",
			in: []string{
				"//kok:op POST /logs",
				`//kok:param ip in=header name=X-Forwarded-For; \`,
				"               in=request name=RemoteAddr",
			},
			want: []string{
				"//kok:op POST /logs",
				"//kok:param ip in=header name=X-Forwarded-For; in=request name=RemoteAddr",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.JoinComments()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Doc: got (%#v), want (%#v)", got, tt.want)
			}
		})
	}
}
