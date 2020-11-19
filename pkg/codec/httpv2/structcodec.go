package codec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagName = "kok"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrMissingRequired = errors.New("missing required field")
)

func getFieldName(field reflect.StructField) (name string, required, omitted bool) {
	kokTag := field.Tag.Get(tagName)
	parts := strings.SplitN(kokTag, ",", 2)

	kokName := parts[0]
	if len(parts) == 2 && parts[1] == "required" {
		required = true
	}

	switch kokName {
	case "":
		name = field.Name
	case "-":
		omitted = true
	default:
		name = kokName
	}

	return
}

// DecodeMapToStruct decodes a value from map[string]string to struct (or *struct).
func DecodeMapToStruct(in map[string]string, out interface{}) error {
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

		fieldName, required, omitted := getFieldName(field)
		if omitted {
			continue
		}

		value := in[fieldName]
		if value == "" {
			if !required {
				continue
			}
			return ErrMissingRequired
		}

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetUint(v)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldValue.SetBool(v)
		case reflect.String:
			fieldValue.SetString(value)
		default:
			panic(fmt.Errorf("unsupported field value: %v", fieldValue))
		}
	}

	return nil
}

// DecodeMapToStruct encode a value from struct (or *struct) to map[string]string.
func EncodeStructToMap(in interface{}, out *map[string]string) error {
	inValue := reflect.ValueOf(in)
	switch k := inValue.Kind(); {
	case k == reflect.Ptr && inValue.Elem().Kind() == reflect.Struct:
		// Convert inValue from *struct to struct implicitly.
		inValue = inValue.Elem()
	case k == reflect.Struct:
	default:
		return ErrUnsupportedType
	}

	if out == nil {
		panic(fmt.Errorf("invalid out: %#v", out))
	}
	outMap := *out

	structType := inValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := inValue.Field(i)

		fieldName, _, omitted := getFieldName(field)
		if omitted {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v := fieldValue.Int()
			outMap[fieldName] = strconv.FormatInt(v, 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v := fieldValue.Uint()
			outMap[fieldName] = strconv.FormatUint(v, 10)
		case reflect.Bool:
			v := fieldValue.Bool()
			outMap[fieldName] = strconv.FormatBool(v)
		case reflect.String:
			outMap[fieldName] = fieldValue.String()
		default:
			panic(fmt.Errorf("unsupported field value: %v", fieldValue))
		}
	}

	return nil
}
