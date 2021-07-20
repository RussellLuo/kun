package moq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	cases := []struct {
		SrcDir    string
		IfaceName string
	}{
		{
			SrcDir:    "./testdata/simple",
			IfaceName: "Counter",
		},
		{
			SrcDir:    "./testdata/importalias",
			IfaceName: "Connector",
		},
	}

	for _, c := range cases {
		parser, err := New(Config{SrcDir: c.SrcDir})
		if err != nil {
			t.Fatalf("err: %v\n", err)
		}

		data, err := parser.Parse(c.IfaceName)
		if err != nil {
			t.Fatalf("err: %v\n", err)
		}

		gotJSON, err := json.MarshalIndent(data, "", "  ") // 2 spaces indent
		if err != nil {
			t.Fatalf("err: %v\n", err)
		}

		wantJSON, err := ioutil.ReadFile(c.SrcDir + "/data.json")
		if err != nil {
			t.Fatalf("err: %v\n", err)
		}

		if !bytes.Equal(gotJSON, wantJSON) {
			t.Fatalf("Data: got (%s), want (%s)", gotJSON, wantJSON)
		}
	}
}
