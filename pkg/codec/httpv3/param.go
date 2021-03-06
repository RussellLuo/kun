package codec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	tagName = "kok"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrMissingRequired = errors.New("missing required field")
)

// DecodeSliceToBasic decodes a string slice to a basic value (or a slice of basic values).
func DecodeSliceToBasic(values []string, out interface{}) error {
	if len(values) == 0 {
		return nil
	}

	switch v := out.(type) {
	case *int:
		vv, err := strconv.Atoi(values[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]int:
		for _, value := range values {
			vv, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *int8:
		vv, err := strconv.ParseInt(values[0], 10, 8)
		if err != nil {
			return err
		}
		*v = int8(vv)
	case *[]int8:
		for _, value := range values {
			vv, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return err
			}
			*v = append(*v, int8(vv))
		}
	case *int16:
		vv, err := strconv.ParseInt(values[0], 10, 16)
		if err != nil {
			return err
		}
		*v = int16(vv)
	case *[]int16:
		for _, value := range values {
			vv, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}
			*v = append(*v, int16(vv))
		}
	case *int32:
		vv, err := strconv.ParseInt(values[0], 10, 32)
		if err != nil {
			return err
		}
		*v = int32(vv)
	case *[]int32:
		for _, value := range values {
			vv, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			*v = append(*v, int32(vv))
		}
	case *int64:
		vv, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return err
		}
		*v = vv
	case *[]int64:
		for _, value := range values {
			vv, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *uint:
		vv, err := strconv.ParseUint(values[0], 10, 0)
		if err != nil {
			return err
		}
		*v = uint(vv)
	case *[]uint:
		for _, value := range values {
			vv, err := strconv.ParseUint(value, 10, 0)
			if err != nil {
				return err
			}
			*v = append(*v, uint(vv))
		}
	case *uint8:
		vv, err := strconv.ParseUint(values[0], 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(vv)
	case *[]uint8:
		for _, value := range values {
			vv, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return err
			}
			*v = append(*v, uint8(vv))
		}
	case *uint16:
		vv, err := strconv.ParseUint(values[0], 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(vv)
	case *[]uint16:
		for _, value := range values {
			vv, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return err
			}
			*v = append(*v, uint16(vv))
		}
	case *uint32:
		vv, err := strconv.ParseUint(values[0], 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(vv)
	case *[]uint32:
		for _, value := range values {
			vv, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return err
			}
			*v = append(*v, uint32(vv))
		}
	case *uint64:
		vv, err := strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return err
		}
		*v = vv
	case *[]uint64:
		for _, value := range values {
			vv, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *bool:
		vv, err := strconv.ParseBool(values[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]bool:
		for _, value := range values {
			vv, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *string:
		*v = values[0]
	case *[]string:
		*v = append(*v, values...)
	case *time.Time:
		vv, err := time.Parse(time.RFC3339, values[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]time.Time:
		for _, value := range values {
			vv, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *time.Duration:
		vv, err := time.ParseDuration(values[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]time.Duration:
		for _, value := range values {
			vv, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	default:
		// Panic since this is a programming error.
		panic(fmt.Errorf("unsupported out type: %T", v))
	}

	return nil
}

// EncodeBasicToSlice encodes a basic value (or a slice of basic values) to a string slice.
func EncodeBasicToSlice(in interface{}) (values []string) {
	switch v := in.(type) {
	case int:
		values = append(values, strconv.FormatInt(int64(v), 10))
	case []int:
		for _, vv := range v {
			values = append(values, strconv.FormatInt(int64(vv), 10))
		}
	case int8:
		values = append(values, strconv.FormatInt(int64(v), 10))
	case []int8:
		for _, vv := range v {
			values = append(values, strconv.FormatInt(int64(vv), 10))
		}
	case int16:
		values = append(values, strconv.FormatInt(int64(v), 10))
	case []int16:
		for _, vv := range v {
			values = append(values, strconv.FormatInt(int64(vv), 10))
		}
	case int32:
		values = append(values, strconv.FormatInt(int64(v), 10))
	case []int32:
		for _, vv := range v {
			values = append(values, strconv.FormatInt(int64(vv), 10))
		}
	case int64:
		values = append(values, strconv.FormatInt(v, 10))
	case []int64:
		for _, vv := range v {
			values = append(values, strconv.FormatInt(vv, 10))
		}
	case uint:
		values = append(values, strconv.FormatUint(uint64(v), 10))
	case []uint:
		for _, vv := range v {
			values = append(values, strconv.FormatUint(uint64(vv), 10))
		}
	case uint8:
		values = append(values, strconv.FormatUint(uint64(v), 10))
	case []uint8:
		for _, vv := range v {
			values = append(values, strconv.FormatUint(uint64(vv), 10))
		}
	case uint16:
		values = append(values, strconv.FormatUint(uint64(v), 10))
	case []uint16:
		for _, vv := range v {
			values = append(values, strconv.FormatUint(uint64(vv), 10))
		}
	case uint32:
		values = append(values, strconv.FormatUint(uint64(v), 10))
	case []uint32:
		for _, vv := range v {
			values = append(values, strconv.FormatUint(uint64(vv), 10))
		}
	case uint64:
		values = append(values, strconv.FormatUint(v, 10))
	case []uint64:
		for _, vv := range v {
			values = append(values, strconv.FormatUint(vv, 10))
		}
	case bool:
		values = append(values, strconv.FormatBool(v))
	case []bool:
		for _, vv := range v {
			values = append(values, strconv.FormatBool(vv))
		}
	case string:
		values = append(values, v)
	case []string:
		values = v
	case time.Time:
		values = append(values, v.Format(time.RFC3339))
	case []time.Time:
		for _, vv := range v {
			values = append(values, vv.Format(time.RFC3339))
		}
	case time.Duration:
		values = append(values, v.String())
	case []time.Duration:
		for _, vv := range v {
			values = append(values, vv.String())
		}
	default:
		values = append(values, fmt.Sprintf("%v", in))
	}
	return
}

// DecodeMapToStruct decodes a string map to a struct (or a *struct).
func DecodeMapToStruct(in map[string][]string, out interface{}) error {
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

		fieldName, required, omitted := GetFieldName(field)
		if omitted {
			continue
		}

		values := in[fieldName]
		if len(values) == 0 {
			if !required {
				continue
			}
			return ErrMissingRequired
		}

		fieldValuePtr := reflect.New(fieldValue.Type())
		if err := DecodeSliceToBasic(values, fieldValuePtr.Interface()); err != nil {
			return err
		}
		fieldValue.Set(fieldValuePtr.Elem())
	}

	return nil
}

// EncodeStructToMap encodes a struct (or a *struct) to a string map.
func EncodeStructToMap(in interface{}, out *map[string][]string) error {
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

		fieldName, _, omitted := GetFieldName(field)
		if omitted {
			continue
		}

		outMap[fieldName] = EncodeBasicToSlice(fieldValue.Interface())
	}

	return nil
}

func GetFieldName(field reflect.StructField) (name string, required, omitted bool) {
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
