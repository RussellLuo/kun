package httpcodec

import (
	"reflect"
	"strconv"
	"testing"
)

type testOut struct {
	N1 int
	N2 int
	N3 int
}

type plusTen struct{}

func (pt plusTen) Decode(in []string, out interface{}) error {
	n, _ := strconv.Atoi(in[0])

	outPtr := out.(*int)
	*outPtr = n + 10

	return nil
}

func (pt plusTen) Encode(in interface{}) (out []string) {
	v := in.(int) - 10
	return []string{strconv.FormatInt(int64(v), 10)}
}

type eachPlusTen struct{}

func (ept eachPlusTen) Decode(in map[string][]string, out interface{}) error {
	n1, _ := strconv.Atoi(in["query.N1"][0])
	n2, _ := strconv.Atoi(in["query.N2"][0])
	n3, _ := strconv.Atoi(in["query.N3"][0])

	outPtr := out.(*testOut)
	outPtr.N1 = n1 + 10
	outPtr.N2 = n2 + 10
	outPtr.N3 = n3 + 10

	return nil
}

func (ept eachPlusTen) Encode(in interface{}) (out map[string][]string) {
	v := in.(testOut)
	return map[string][]string{
		"query.N1": {strconv.FormatInt(int64(v.N1-10), 10)},
		"query.N2": {strconv.FormatInt(int64(v.N2-10), 10)},
		"query.N3": {strconv.FormatInt(int64(v.N3-10), 10)},
	}
}

func TestPatcher_DecodeRequestParam(t *testing.T) {
	cases := []struct {
		name      string
		inPatcher *Patcher
		inName    string
		inValues  []string
		inOut     int
		wantOut   int
	}{
		{
			name:      "no patching",
			inPatcher: NewPatcher(JSON{}),
			inName:    "arg",
			inValues:  []string{"1"},
			inOut:     0, // an empty int to store the result.
			wantOut:   1,
		},
		{
			name:      "patching param",
			inPatcher: NewPatcher(JSON{}).Param("arg", plusTen{}),
			inName:    "arg",
			inValues:  []string{"1"},
			inOut:     0, // an empty int to store the result.
			wantOut:   11,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.inPatcher.DecodeRequestParam(c.inName, c.inValues, &c.inOut)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if c.inOut != c.wantOut {
				t.Fatalf("Out: got (%#v), want (%#v)", c.inOut, c.wantOut)
			}
		})
	}
}

func TestPatcher_EncodeRequestParam(t *testing.T) {
	cases := []struct {
		name      string
		inPatcher *Patcher
		inName    string
		inValue   interface{}
		wantOut   []string
	}{
		{
			name:      "no patching",
			inPatcher: NewPatcher(JSON{}),
			inName:    "arg",
			inValue:   11,
			wantOut:   []string{"11"},
		},
		{
			name:      "patching param",
			inPatcher: NewPatcher(JSON{}).Param("arg", plusTen{}),
			inName:    "arg",
			inValue:   11,
			wantOut:   []string{"1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.inPatcher.EncodeRequestParam(c.inName, c.inValue)
			if !reflect.DeepEqual(got, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", got, c.wantOut)
			}
		})
	}
}

func TestPatcher_DecodeRequestParams(t *testing.T) {
	cases := []struct {
		name      string
		inPatcher *Patcher
		inName    string
		inValues  map[string][]string
		inOut     testOut
		wantOut   testOut
	}{
		{
			name:      "no patching",
			inPatcher: NewPatcher(JSON{}),
			inName:    "args",
			inValues:  map[string][]string{"query.N1": {"1"}, "query.N2": {"2"}, "query.N3": {"3"}},
			inOut:     testOut{},
			wantOut:   testOut{N1: 1, N2: 2, N3: 3},
		},
		{
			name:      "patching params",
			inPatcher: NewPatcher(JSON{}).Params("args", eachPlusTen{}),
			inName:    "args",
			inValues:  map[string][]string{"query.N1": {"1"}, "query.N2": {"2"}, "query.N3": {"3"}},
			inOut:     testOut{},
			wantOut:   testOut{N1: 11, N2: 12, N3: 13},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.inPatcher.DecodeRequestParams(c.inName, c.inValues, &c.inOut)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			if c.inOut != c.wantOut {
				t.Fatalf("Out: got (%#v), want (%#v)", c.inOut, c.wantOut)
			}
		})
	}
}

func TestPatcher_EncodeRequestParams(t *testing.T) {
	cases := []struct {
		name      string
		inPatcher *Patcher
		inName    string
		inValue   interface{}
		wantOut   map[string][]string
	}{
		{
			name:      "no patching",
			inPatcher: NewPatcher(JSON{}),
			inName:    "args",
			inValue:   testOut{N1: 11, N2: 12, N3: 13},
			wantOut:   map[string][]string{"query.N1": {"11"}, "query.N2": {"12"}, "query.N3": {"13"}},
		},
		{
			name:      "patching params",
			inPatcher: NewPatcher(JSON{}).Params("args", eachPlusTen{}),
			inName:    "args",
			inValue:   testOut{N1: 11, N2: 12, N3: 13},
			wantOut:   map[string][]string{"query.N1": {"1"}, "query.N2": {"2"}, "query.N3": {"3"}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.inPatcher.EncodeRequestParams(c.inName, c.inValue)
			if !reflect.DeepEqual(got, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", got, c.wantOut)
			}
		})
	}
}
