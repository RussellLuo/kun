package httpcodec

import (
	"reflect"
	"testing"
	"time"
)

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
		"query.int":      {"1"},
		"query.ints":     {"1", "2"},
		"query.int8":     {"2"},
		"query.int8s":    {"2", "3"},
		"query.int16":    {"3"},
		"query.int16s":   {"3", "4"},
		"query.int32":    {"4"},
		"query.int32s":   {"4", "5"},
		"query.int64":    {"5"},
		"query.int64s":   {"5", "6"},
		"query.uint":     {"6"},
		"query.uints":    {"6", "7"},
		"query.uint8":    {"7"},
		"query.uint8s":   {"7", "8"},
		"query.uint16":   {"8"},
		"query.uint16s":  {"8", "9"},
		"query.uint32":   {"9"},
		"query.uint32s":  {"9", "10"},
		"query.uint64":   {"10"},
		"query.uint64s":  {"10", "11"},
		"query.bool":     {"true"},
		"query.bools":    {"true", "false"},
		"query.string":   {"hello"},
		"query.strings":  {"hello", "hi"},
		"query.required": {"wow"},
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
				"query.string":   nil,
				"query.required": {"wow"},
			},
			outPtr:  ptrToValue,
			wantOut: value{Required: "wow"},
		},
		{
			name: "missing required field",
			in: map[string][]string{
				"query.required": nil,
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
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
			},
			outPtr:  nil,
			wantErr: ErrUnsupportedType,
		},
		{
			name: "struct",
			in: map[string][]string{
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
			},
			outPtr:  value{},
			wantErr: ErrUnsupportedType,
		},
		{
			name: "string",
			in: map[string][]string{
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
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
		"query.int":     {"1"},
		"query.ints":    {"1", "2"},
		"query.int8":    {"2"},
		"query.int8s":   {"2", "3"},
		"query.int16":   {"3"},
		"query.int16s":  {"3", "4"},
		"query.int32":   {"4"},
		"query.int32s":  {"4", "5"},
		"query.int64":   {"5"},
		"query.int64s":  {"5", "6"},
		"query.uint":    {"6"},
		"query.uints":   {"6", "7"},
		"query.uint8":   {"7"},
		"query.uint8s":  {"7", "8"},
		"query.uint16":  {"8"},
		"query.uint16s": {"8", "9"},
		"query.uint32":  {"9"},
		"query.uint32s": {"9", "10"},
		"query.uint64":  {"10"},
		"query.uint64s": {"10", "11"},
		"query.bool":    {"true"},
		"query.bools":   {"true", "false"},
		"query.string":  {"hello"},
		"query.strings": {"hello", "hi"},
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

func TestGetKokField(t *testing.T) {
	cases := []struct {
		name    string
		in      reflect.StructField
		wantOut KokField
	}{
		{
			name:    "in path",
			in:      reflect.StructField{Name: "ID", Tag: `kok:"path.id"`},
			wantOut: KokField{Name: "path.id", Required: true},
		},
		{
			name:    "in query",
			in:      reflect.StructField{Name: "ID", Tag: `kok:"query.id"`},
			wantOut: KokField{Name: "query.id"},
		},
		{
			name:    "omitted",
			in:      reflect.StructField{Name: "ID", Tag: `kok:"-"`},
			wantOut: KokField{Omitted: true},
		},
		{
			name:    "required",
			in:      reflect.StructField{Name: "ID", Tag: `kok:",required"`},
			wantOut: KokField{Name: "query.ID", Required: true},
		},
		{
			name:    "has type",
			in:      reflect.StructField{Name: "ID", Tag: `kok:",type:string"`},
			wantOut: KokField{Name: "query.ID", Type: "string"},
		},
		{
			name:    "has description",
			in:      reflect.StructField{Name: "ID", Tag: `kok:",descr:string"`},
			wantOut: KokField{Name: "query.ID", Description: "string"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := GetKokField(c.in)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}

func TestBasicParam_Decode(t *testing.T) {
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
		Float32   float32
		Float32s  []float32
		Float64   float64
		Float64s  []float64
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
			name:    "float32",
			values:  []string{"1", "2"},
			out:     &v.Float32,
			wantOut: float32(1),
		},
		{
			name:    "[]float32",
			values:  []string{"1", "2"},
			out:     &v.Float32s,
			wantOut: []float32{1, 2},
		},
		{
			name:    "float64",
			values:  []string{"1", "2"},
			out:     &v.Float64,
			wantOut: float64(1),
		},
		{
			name:    "[]float64",
			values:  []string{"1", "2"},
			out:     &v.Float64s,
			wantOut: []float64{1, 2},
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

	p := BasicParam{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := p.Decode(c.values, c.out)
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

func TestBasicParam_Encode(t *testing.T) {
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
			name:    "float32",
			value:   float32(1),
			wantOut: []string{"1"},
		},
		{
			name:    "[]float32",
			value:   []float32{1, 2},
			wantOut: []string{"1", "2"},
		},
		{
			name:    "float64",
			value:   float64(1),
			wantOut: []string{"1"},
		},
		{
			name:    "[]float64",
			value:   []float64{1, 2},
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

	p := BasicParam{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := p.Encode(c.value)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}

func TestStructParams_Decode(t *testing.T) {
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
		"query.int":      {"1"},
		"query.ints":     {"1", "2"},
		"query.int8":     {"2"},
		"query.int8s":    {"2", "3"},
		"query.int16":    {"3"},
		"query.int16s":   {"3", "4"},
		"query.int32":    {"4"},
		"query.int32s":   {"4", "5"},
		"query.int64":    {"5"},
		"query.int64s":   {"5", "6"},
		"query.uint":     {"6"},
		"query.uints":    {"6", "7"},
		"query.uint8":    {"7"},
		"query.uint8s":   {"7", "8"},
		"query.uint16":   {"8"},
		"query.uint16s":  {"8", "9"},
		"query.uint32":   {"9"},
		"query.uint32s":  {"9", "10"},
		"query.uint64":   {"10"},
		"query.uint64s":  {"10", "11"},
		"query.bool":     {"true"},
		"query.bools":    {"true", "false"},
		"query.string":   {"hello"},
		"query.strings":  {"hello", "hi"},
		"query.required": {"wow"},
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
				"query.string":   nil,
				"query.required": {"wow"},
			},
			outPtr:  ptrToValue,
			wantOut: value{Required: "wow"},
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
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
			},
			outPtr:  nil,
			wantErr: ErrUnsupportedType,
		},
		{
			name: "struct",
			in: map[string][]string{
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
			},
			outPtr:  value{},
			wantErr: ErrUnsupportedType,
		},
		{
			name: "string",
			in: map[string][]string{
				"query.int":    {"1"},
				"query.uint":   {"6"},
				"query.bool":   {"true"},
				"query.string": {"hello"},
			},
			outPtr:  new(string),
			wantErr: ErrUnsupportedType,
		},
	}

	p := StructParams{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := p.Decode(c.in, c.outPtr)
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

func TestStructParams_Encode(t *testing.T) {
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
		"query.int":     {"1"},
		"query.ints":    {"1", "2"},
		"query.int8":    {"2"},
		"query.int8s":   {"2", "3"},
		"query.int16":   {"3"},
		"query.int16s":  {"3", "4"},
		"query.int32":   {"4"},
		"query.int32s":  {"4", "5"},
		"query.int64":   {"5"},
		"query.int64s":  {"5", "6"},
		"query.uint":    {"6"},
		"query.uints":   {"6", "7"},
		"query.uint8":   {"7"},
		"query.uint8s":  {"7", "8"},
		"query.uint16":  {"8"},
		"query.uint16s": {"8", "9"},
		"query.uint32":  {"9"},
		"query.uint32s": {"9", "10"},
		"query.uint64":  {"10"},
		"query.uint64s": {"10", "11"},
		"query.bool":    {"true"},
		"query.bools":   {"true", "false"},
		"query.string":  {"hello"},
		"query.strings": {"hello", "hi"},
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
	}

	p := StructParams{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := p.Encode(c.in)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}
