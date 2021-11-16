package httpcodec

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"testing"
)

func TestDecodeMultipartFormToStruct(t *testing.T) {
	fromString := func(name, value string) *FormFile {
		return &FormFile{
			Name: name,
			Size: int64(len(value)),
			File: ioutil.NopCloser(bytes.NewBufferString(value)),
		}
	}

	type form struct {
		Text string    `json:"text"`
		File *FormFile `json:"file"`
	}

	inFileContent := "this is a file"
	in := form{
		Text: "this is a text",
		File: fromString("filename", inFileContent),
	}

	// 1. Encode the struct named in to a multipart message.
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	if err := encodeStructToMultipartForm(in, writer); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// 2. Decode the multipart message to a struct named out.

	// Parse the boundary.
	_, params, err := mime.ParseMediaType(writer.FormDataContentType())
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// Parse the multipart message.
	reader := multipart.NewReader(data, params["boundary"])
	formData, err := reader.ReadForm(defaultMaxMemory)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	out := &form{}
	if err := decodeMultipartFormToStruct(formData, out); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// 3. Assert the equality of out and in.
	if out.Text != in.Text {
		t.Fatalf("Name: got (%#v), want (%#v)", out.Text, in.Text)
	}
	if out.File.Name != in.File.Name {
		t.Fatalf("File.Name: got (%#v), want (%#v)", out.File.Name, in.File.Name)
	}
	outFileBytes, _ := ioutil.ReadAll(out.File.File)
	if string(outFileBytes) != inFileContent {
		t.Fatalf("File.File: got (%s), want (%s)", outFileBytes, inFileContent)
	}
}
