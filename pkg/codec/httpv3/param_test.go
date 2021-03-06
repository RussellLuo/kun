package codec

import (
	"reflect"
	"testing"
	"time"
)

func Test_Decode(t *testing.T) {
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
			err := Decode(c.values, c.out)
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

func Test_Encode(t *testing.T) {
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
			out := Encode(c.value)
			if !reflect.DeepEqual(out, c.wantOut) {
				t.Fatalf("Out: got (%#v), want (%#v)", out, c.wantOut)
			}
		})
	}
}
