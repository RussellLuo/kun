package openapi

import (
	"reflect"
	"testing"
)

func Test_extractPathVarNames(t *testing.T) {
	want := []string{"name", "path"}
	got := extractPathVarNames("/var/{name}/in/{path}")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Names: got (%#v), want (%#v)", got, want)
	}
}
