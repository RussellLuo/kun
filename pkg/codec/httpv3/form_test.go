package codec

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
			File: ioutil.NopCloser(bytes.NewBufferString(value)),
		}
	}

	type form struct {
		Name string    `json:"name"`
		File *FormFile `json:"file"`
	}

	inFileContent := "value string"
	in := form{
		Name: "file",
		File: fromString("filename", inFileContent),
	}

	// 1. Encode the struct in to a multipart message.
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	if err := encodeStructToMultipartForm(in, writer); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// 2. Decode the multipart message to a struct out.

	// Create an instance of multipart.Form.
	_, params, err := mime.ParseMediaType(writer.FormDataContentType())
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	reader := multipart.NewReader(data, params["boundary"])
	formData, err := reader.ReadForm(defaultMaxMemory)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	out := &form{}
	if err := decodeMultipartFormToStruct(formData, out); err != nil {
		t.Fatalf("Error: %v", err)
	}

	// 3. Validate the equality of out and in.
	if out.Name != in.Name {
		t.Fatalf("Name: got (%#v), want (%#v)", out.Name, in.Name)
	}
	if out.File.Name != in.File.Name {
		t.Fatalf("File.Name: got (%#v), want (%#v)", out.File.Name, in.File.Name)
	}
	outFileBytes, _ := ioutil.ReadAll(out.File.File)
	if string(outFileBytes) != inFileContent {
		t.Fatalf("File.File: got (%s), want (%s)", outFileBytes, inFileContent)
	}
}
