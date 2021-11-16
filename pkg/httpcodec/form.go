package httpcodec

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// TODO: add support for PostForm

type MultipartForm struct {
	JSON

	maxMemory int64
}

func NewMultipartForm(maxMemory int64) *MultipartForm {
	if maxMemory == 0 {
		maxMemory = defaultMaxMemory
	}
	return &MultipartForm{
		maxMemory: maxMemory,
	}
}

func (mf *MultipartForm) DecodeRequestBody(r *http.Request, out interface{}) error {
	if err := r.ParseMultipartForm(mf.maxMemory); err != nil {
		return err
	}
	return decodeMultipartFormToStruct(r.MultipartForm, out)
}

func (mf *MultipartForm) EncodeRequestBody(in interface{}) (io.Reader, map[string]string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := encodeStructToMultipartForm(in, writer); err != nil {
		return nil, nil, err
	}
	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	return body, headers, nil
}

// decodeMultipartFormToStruct decodes a multipart message to a struct (or a *struct).
func decodeMultipartFormToStruct(form *multipart.Form, out interface{}) error {
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr || outValue.IsNil() {
		return ErrUnsupportedType
	}

	elemValue := outValue.Elem()
	elemType := elemValue.Type()

	var structValue reflect.Value

	switch k := elemValue.Kind(); {
	case k == reflect.Struct:
		structValue = elemValue
	case k == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct:
		// To handle possible nil pointer, always create a pointer
		// to a new zero struct.
		structValuePtr := reflect.New(elemType.Elem())
		outValue.Elem().Set(structValuePtr)

		structValue = structValuePtr.Elem()
	default:
		return ErrUnsupportedType
	}

	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		fieldName, omitted := getFormFieldName(field)
		if omitted {
			continue
		}

		fieldValuePtr := reflect.New(fieldValue.Type())

		switch v := fieldValuePtr.Interface().(type) {
		case **FormFile:
			// Decode the first form file.
			if files, ok := form.File[fieldName]; ok && len(files) > 0 {
				f, err := FromMultipartFileHeader(files[0])
				if err != nil {
					return err
				}
				*v = f
			}

		case *[]*FormFile:
			// Decode all the form files.
			if files, ok := form.File[fieldName]; ok {
				for _, file := range files {
					f, err := FromMultipartFileHeader(file)
					if err != nil {
						return err
					}
					*v = append(*v, f)
				}
			}

		default:
			// Decode normal form values.
			values := form.Value[fieldName]
			if err := defaultBasicParam.Decode(values, fieldValuePtr.Interface()); err != nil {
				return err
			}
		}

		fieldValue.Set(fieldValuePtr.Elem())
	}

	return nil
}

// encodeStructToMultipartForm encodes a struct (or a *struct) to a multipart message.
func encodeStructToMultipartForm(in interface{}, writer *multipart.Writer) error {
	inValue := reflect.ValueOf(in)
	switch k := inValue.Kind(); {
	case k == reflect.Ptr && inValue.Elem().Kind() == reflect.Struct:
		// Convert inValue from *struct to struct implicitly.
		inValue = inValue.Elem()
	case k == reflect.Struct:
	default:
		return ErrUnsupportedType
	}

	if writer == nil {
		return errors.New("writer is nil")
	}

	structType := inValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := inValue.Field(i)

		fieldName, omitted := getFormFieldName(field)
		if omitted {
			continue
		}

		switch v := fieldValue.Interface().(type) {
		case *FormFile:
			// Write the first form file.
			if err := writeFile(writer, fieldName, v); err != nil {
				return err
			}
		case []*FormFile:
			// Write all the form files.
			for _, vv := range v {
				if err := writeFile(writer, fieldName, vv); err != nil {
					return err
				}
			}
		default:
			// Write normal form values.
			for _, value := range defaultBasicParam.Encode(fieldValue.Interface()) {
				if err := writer.WriteField(fieldName, value); err != nil {
					return err
				}
			}
		}
	}

	// Finish the multipart message.
	return writer.Close()
}

func writeFile(writer *multipart.Writer, fieldName string, file *FormFile) error {
	fileWriter, err := writer.CreateFormFile(fieldName, file.Name)
	if err != nil {
		return err
	}

	defer file.File.Close()
	if _, err := io.Copy(fileWriter, file.File); err != nil {
		return err
	}
	return nil
}

func getFormFieldName(field reflect.StructField) (name string, omitted bool) {
	kokTag := field.Tag.Get("json")
	parts := strings.SplitN(kokTag, ",", 2)

	name = parts[0]

	switch name {
	case "":
		name = field.Name
	case "-":
		omitted = true
	}

	return
}

// FormFile describes a file part of a multipart message.
type FormFile struct {
	Name   string
	Header map[string][]string
	Size   int64
	File   io.ReadCloser
}

func FromMultipartFileHeader(fh *multipart.FileHeader) (*FormFile, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	return &FormFile{
		Name:   fh.Filename,
		Header: fh.Header,
		Size:   fh.Size,
		File:   file,
	}, nil
}

func FromOSFile(name string) (*FormFile, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()

	return &FormFile{
		Name: name,
		// Header is nil for OS files.
		Size: size,
		File: file,
	}, nil
}

func (ff *FormFile) Save(name string) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	defer ff.File.Close()
	if _, err := io.Copy(out, ff.File); err != nil {
		return err
	}
	return nil
}
