package httpcodec

import (
	"reflect"
	"testing"
	"time"
)

func TestDecodeSliceToBasic(t *testing.T) {
	type value struct {
		Int       int
		Ints      []int
		Int8      int8
		Int8s     []int8
		Int16     int16
		Int16s    []int16
		Int32     int32
		Int32s    []int32
		Int64     int64
		Int64s    []int64
		Uint      uint
		Uints     []uint
		Uint8     uint8
		Uint8s    []uint8
		Uint16    uint16
		Uint16s   []uint16
		Uint32    uint32
		Uint32s   []uint32
		Uint64    uint64
		Uint64s   []uint64
		Bool      bool
		Bools     []bool
		String    string
		Strings   []string
		Time      time.Time
		Times     []time.Time
		Duration  time.Duration
		Durations []time.Duration
	}
	v := value{}

	cases := []struct {
		name    string
		values  []string
		out     interface{}
		wantOut interface{}
		wantErr error
	}{
		{
			name:    "int",
			values:  []string{"1"},
			out:     &v.Int,
			wantOut: 1,
		},
		{
			name:    "[]int",
			values:  []string{"1", "2"},
			out:     &v.Ints,
			wantOut: []int{1, 2},
		},
		{
			name:    "uint",
			values:  []string{"1"},
			out:     &v.Uint,
			wantOut: uint(1),
		},
		{
			name:    "[]uint",
			values:  []string{"1", "2"},
			out:     &v.Uints,
			wantOut: []uint{1, 2},
		},
		{
			name:    "bool",
			values:  []string{"true"},
			out:     &v.Bool,
			wantOut: true,
		},
		{
			name:    "[]bool",
			values:  []string{"true", "false"},
			out:     &v.Bools,
			wantOut: []bool{true, false},
		},
		{
			name:    "string",
			values:  []string{"yes"},
			out:     &v.String,
			wantOut: "yes",
		},
		{
			name:    "CSV string",
			values:  []string{"v0=0,v1=1"},
			out:     &v.String,
			wantOut: "v0=0,v1=1",
		},
		{
			name:    "empty value for string",
			values:  []string{""},
			out:     &v.String,
			wantOut: "",
		},
		{
			name:    "[]string",
			values:  []string{"yes", "no"},
			out:     &v.Strings,
			wantOut: []string{"yes", "no"},
		},
		{
			name:    "duration",
			values:  []string{"2s"},
			out:     &v.Duration,
			wantOut: 2 * time.Second,
		},
		{
			name:    "[]duration",
			values:  []string{"2s", "4m"},
			out:     &v.Durations,
			wantOut: []time.Duration{2 * time.Second, 4 * time.Minute},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := DecodeSliceToBasic(c.values, c.out)
			if err != c.wantErr {
				t.Fatalf("Err: got (%#v), want (%#v)", err, c.wantErr)
			}
			if err == nil {
				out := reflect.ValueOf(c.out).Elem().Interface()
				if !reflect.DeepEqual(out, c.wantOut) {
					t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
				}
			}
		})
	}
}

func TestEncodeBasicToSlice(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		wantOut []string
	}{
		{
			name:    "int",
			value:   1,
			wantOut: []string{"1"},
		},
		{
			name:    "[]int",
			value:   []int{1, 2},
			wantOut: []string{"1", "2"},
		},
		{
			name:    "uint",
			value:   uint(1),
			wantOut: []string{"1"},
		},
		{
			name:    "[]uint",
			value:   []uint{1, 2},
			wantOut: []string{"1", "2"},
		},
		{
			name:    "bool",
			value:   true,
			wantOut: []string{"true"},
		},
		{
			name:    "[]bool",
			value:   []bool{true, false},
			wantOut: []string{"true", "false"},
		},
		{
			name:    "string",
			value:   "yes",
			wantOut: []string{"yes"},
		},
		{
			name:    "[]string",
			value:   []string{"yes", "no"},
			wantOut: []string{"yes", "no"},
		},
		{
			name:    "duration",
			value:   2 * time.Second,
			wantOut: []string{"2s"},
		},
		{
			name:    "[]duration",
			value:   []time.Duration{2 * time.Second, 4 * time.Minute},
			wantOut: []string{"2s", "4m0s"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := EncodeBasicToSlice(c.value)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}

func TestDecodeMapToStruct(t *testing.T) {
	type value struct {
		Int      int      `kok:"int"`
		Ints     []int    `kok:"ints"`
		Int8     int8     `kok:"int8"`
		Int8s    []int8   `kok:"int8s"`
		Int16    int16    `kok:"int16"`
		Int16s   []int16  `kok:"int16s"`
		Int32    int32    `kok:"int32"`
		Int32s   []int32  `kok:"int32s"`
		Int64    int64    `kok:"int64"`
		Int64s   []int64  `kok:"int64s"`
		Uint     uint     `kok:"uint"`
		Uints    []uint   `kok:"uints"`
		Uint8    uint8    `kok:"uint8"`
		Uint8s   []uint8  `kok:"uint8s"`
		Uint16   uint16   `kok:"uint16"`
		Uint16s  []uint16 `kok:"uint16s"`
		Uint32   uint32   `kok:"uint32"`
		Uint32s  []uint32 `kok:"uint32s"`
		Uint64   uint64   `kok:"uint64"`
		Uint64s  []uint64 `kok:"uint64s"`
		Bool     bool     `kok:"bool"`
		Bools    []bool   `kok:"bools"`
		String   string   `kok:"string"`
		Strings  []string `kok:"strings"`
		Required string   `kok:"required,required"`
	}

	testIn := map[string][]string{
		"int":      {"1"},
		"ints":     {"1", "2"},
		"int8":     {"2"},
		"int8s":    {"2", "3"},
		"int16":    {"3"},
		"int16s":   {"3", "4"},
		"int32":    {"4"},
		"int32s":   {"4", "5"},
		"int64":    {"5"},
		"int64s":   {"5", "6"},
		"uint":     {"6"},
		"uints":    {"6", "7"},
		"uint8":    {"7"},
		"uint8s":   {"7", "8"},
		"uint16":   {"8"},
		"uint16s":  {"8", "9"},
		"uint32":   {"9"},
		"uint32s":  {"9", "10"},
		"uint64":   {"10"},
		"uint64s":  {"10", "11"},
		"bool":     {"true"},
		"bools":    {"true", "false"},
		"string":   {"hello"},
		"strings":  {"hello", "hi"},
		"required": {"wow"},
	}
	testValue := value{
		Int:      1,
		Ints:     []int{1, 2},
		Int8:     2,
		Int8s:    []int8{2, 3},
		Int16:    3,
		Int16s:   []int16{3, 4},
		Int32:    4,
		Int32s:   []int32{4, 5},
		Int64:    5,
		Int64s:   []int64{5, 6},
		Uint:     6,
		Uints:    []uint{6, 7},
		Uint8:    7,
		Uint8s:   []uint8{7, 8},
		Uint16:   8,
		Uint16s:  []uint16{8, 9},
		Uint32:   9,
		Uint32s:  []uint32{9, 10},
		Uint64:   10,
		Uint64s:  []uint64{10, 11},
		Bool:     true,
		Bools:    []bool{true, false},
		String:   "hello",
		Strings:  []string{"hello", "hi"},
		Required: "wow",
	}

	ptrToValue := &value{}
	var nilPtrToValue *value = nil

	cases := []struct {
		name    string
		in      map[string][]string
		outPtr  interface{}
		wantOut interface{}
		wantErr error
	}{
		{
			name: "missing optional field",
			in: map[string][]string{
				"string":   nil,
				"required": {"wow"},
			},
			outPtr:  ptrToValue,
			wantOut: value{Required: "wow"},
		},
		{
			name: "missing required field",
			in: map[string][]string{
				"required": nil,
			},
			outPtr:  ptrToValue,
			wantErr: ErrMissingRequired,
		},
		{
			name:    "struct pointer",
			in:      testIn,
			outPtr:  ptrToValue,
			wantOut: testValue,
		},
		{
			name:    "pointer of struct pointer",
			in:      testIn,
			outPtr:  &ptrToValue,
			wantOut: &testValue,
		},
		{
			name:    "pointer of nil struct pointer",
			in:      testIn,
			outPtr:  &nilPtrToValue,
			wantOut: &testValue,
		},
		{
			name: "nil",
			in: map[string][]string{
				"int":    {"1"},
				"uint":   {"6"},
				"bool":   {"true"},
				"string": {"hello"},
			},
			outPtr:  nil,
			wantErr: ErrUnsupportedType,
		},
		{
			name: "struct",
			in: map[string][]string{
				"int":    {"1"},
				"uint":   {"6"},
				"bool":   {"true"},
				"string": {"hello"},
			},
			outPtr:  value{},
			wantErr: ErrUnsupportedType,
		},
		{
			name: "string",
			in: map[string][]string{
				"int":    {"1"},
				"uint":   {"6"},
				"bool":   {"true"},
				"string": {"hello"},
			},
			outPtr:  new(string),
			wantErr: ErrUnsupportedType,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := DecodeMapToStruct(c.in, c.outPtr)
			if err != c.wantErr {
				t.Fatalf("Err: got (%#v), want (%#v)", err, c.wantErr)
			}
			if err == nil {
				out := reflect.ValueOf(c.outPtr).Elem().Interface()
				if !reflect.DeepEqual(out, c.wantOut) {
					t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
				}
			}
		})
	}
}

func TestEncodeStructToMap(t *testing.T) {
	type value struct {
		Int     int      `kok:"int"`
		Ints    []int    `kok:"ints"`
		Int8    int8     `kok:"int8"`
		Int8s   []int8   `kok:"int8s"`
		Int16   int16    `kok:"int16"`
		Int16s  []int16  `kok:"int16s"`
		Int32   int32    `kok:"int32"`
		Int32s  []int32  `kok:"int32s"`
		Int64   int64    `kok:"int64"`
		Int64s  []int64  `kok:"int64s"`
		Uint    uint     `kok:"uint"`
		Uints   []uint   `kok:"uints"`
		Uint8   uint8    `kok:"uint8"`
		Uint8s  []uint8  `kok:"uint8s"`
		Uint16  uint16   `kok:"uint16"`
		Uint16s []uint16 `kok:"uint16s"`
		Uint32  uint32   `kok:"uint32"`
		Uint32s []uint32 `kok:"uint32s"`
		Uint64  uint64   `kok:"uint64"`
		Uint64s []uint64 `kok:"uint64s"`
		Bool    bool     `kok:"bool"`
		Bools   []bool   `kok:"bools"`
		String  string   `kok:"string"`
		Strings []string `kok:"strings"`
	}

	testIn := value{
		Int:     1,
		Ints:    []int{1, 2},
		Int8:    2,
		Int8s:   []int8{2, 3},
		Int16:   3,
		Int16s:  []int16{3, 4},
		Int32:   4,
		Int32s:  []int32{4, 5},
		Int64:   5,
		Int64s:  []int64{5, 6},
		Uint:    6,
		Uints:   []uint{6, 7},
		Uint8:   7,
		Uint8s:  []uint8{7, 8},
		Uint16:  8,
		Uint16s: []uint16{8, 9},
		Uint32:  9,
		Uint32s: []uint32{9, 10},
		Uint64:  10,
		Uint64s: []uint64{10, 11},
		Bool:    true,
		Bools:   []bool{true, false},
		String:  "hello",
		Strings: []string{"hello", "hi"},
	}
	testOut := map[string][]string{
		"int":     {"1"},
		"ints":    {"1", "2"},
		"int8":    {"2"},
		"int8s":   {"2", "3"},
		"int16":   {"3"},
		"int16s":  {"3", "4"},
		"int32":   {"4"},
		"int32s":  {"4", "5"},
		"int64":   {"5"},
		"int64s":  {"5", "6"},
		"uint":    {"6"},
		"uints":   {"6", "7"},
		"uint8":   {"7"},
		"uint8s":  {"7", "8"},
		"uint16":  {"8"},
		"uint16s": {"8", "9"},
		"uint32":  {"9"},
		"uint32s": {"9", "10"},
		"uint64":  {"10"},
		"uint64s": {"10", "11"},
		"bool":    {"true"},
		"bools":   {"true", "false"},
		"string":  {"hello"},
		"strings": {"hello", "hi"},
	}

	cases := []struct {
		name    string
		in      interface{}
		wantOut map[string][]string
		wantErr error
	}{
		{
			name:    "struct pointer",
			in:      &testIn,
			wantOut: testOut,
		},
		{
			name:    "struct",
			in:      testIn,
			wantOut: testOut,
		},
		{
			name:    "string",
			in:      "",
			wantErr: ErrUnsupportedType,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := make(map[string][]string)
			err := EncodeStructToMap(c.in, &out)
			if err != c.wantErr {
				t.Fatalf("Err: got (%#v), want (%#v)", err, c.wantErr)
			}
			if err == nil && !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}
