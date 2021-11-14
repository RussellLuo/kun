package httpcodec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/RussellLuo/kok/gen/http/parser"
)

const (
	specialKey = "_httpcodec_special_key"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")

	defaultBasicParam   = BasicParam{}
	defaultBasicParams  = ToParamsCodec(defaultBasicParam)
	defaultStructParams = StructParams{}
)

// BasicParam is a built-in implementation of ParamCodec. It is mainly used
// to encode and decode a basic value or a slice of basic values.
type BasicParam struct{}

// Decode decodes a string slice to a basic value (or a slice of basic values).
func (p BasicParam) Decode(in []string, out interface{}) error {
	if len(in) == 0 {
		return nil
	}

	switch v := out.(type) {
	case *int:
		vv, err := strconv.Atoi(in[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]int:
		for _, value := range in {
			vv, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *int8:
		vv, err := strconv.ParseInt(in[0], 10, 8)
		if err != nil {
			return err
		}
		*v = int8(vv)
	case *[]int8:
		for _, value := range in {
			vv, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return err
			}
			*v = append(*v, int8(vv))
		}
	case *int16:
		vv, err := strconv.ParseInt(in[0], 10, 16)
		if err != nil {
			return err
		}
		*v = int16(vv)
	case *[]int16:
		for _, value := range in {
			vv, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}
			*v = append(*v, int16(vv))
		}
	case *int32:
		vv, err := strconv.ParseInt(in[0], 10, 32)
		if err != nil {
			return err
		}
		*v = int32(vv)
	case *[]int32:
		for _, value := range in {
			vv, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			*v = append(*v, int32(vv))
		}
	case *int64:
		vv, err := strconv.ParseInt(in[0], 10, 64)
		if err != nil {
			return err
		}
		*v = vv
	case *[]int64:
		for _, value := range in {
			vv, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *uint:
		vv, err := strconv.ParseUint(in[0], 10, 0)
		if err != nil {
			return err
		}
		*v = uint(vv)
	case *[]uint:
		for _, value := range in {
			vv, err := strconv.ParseUint(value, 10, 0)
			if err != nil {
				return err
			}
			*v = append(*v, uint(vv))
		}
	case *uint8:
		vv, err := strconv.ParseUint(in[0], 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(vv)
	case *[]uint8:
		for _, value := range in {
			vv, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return err
			}
			*v = append(*v, uint8(vv))
		}
	case *uint16:
		vv, err := strconv.ParseUint(in[0], 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(vv)
	case *[]uint16:
		for _, value := range in {
			vv, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return err
			}
			*v = append(*v, uint16(vv))
		}
	case *uint32:
		vv, err := strconv.ParseUint(in[0], 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(vv)
	case *[]uint32:
		for _, value := range in {
			vv, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return err
			}
			*v = append(*v, uint32(vv))
		}
	case *uint64:
		vv, err := strconv.ParseUint(in[0], 10, 64)
		if err != nil {
			return err
		}
		*v = vv
	case *[]uint64:
		for _, value := range in {
			vv, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *float32:
		vv, err := strconv.ParseFloat(in[0], 32)
		if err != nil {
			return err
		}
		*v = float32(vv)
	case *[]float32:
		for _, value := range in {
			vv, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}
			*v = append(*v, float32(vv))
		}
	case *float64:
		vv, err := strconv.ParseFloat(in[0], 64)
		if err != nil {
			return err
		}
		*v = vv
	case *[]float64:
		for _, value := range in {
			vv, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *bool:
		vv, err := strconv.ParseBool(in[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]bool:
		for _, value := range in {
			vv, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *string:
		*v = in[0]
	case *[]string:
		*v = append(*v, in...)
	case *time.Time:
		vv, err := time.Parse(time.RFC3339, in[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]time.Time:
		for _, value := range in {
			vv, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	case *time.Duration:
		vv, err := time.ParseDuration(in[0])
		if err != nil {
			return err
		}
		*v = vv
	case *[]time.Duration:
		for _, value := range in {
			vv, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			*v = append(*v, vv)
		}
	default:
		return ErrUnsupportedType
	}

	return nil
}

// Encode encodes a basic value (or a slice of basic values) to a string slice.
func (p BasicParam) Encode(in interface{}) (out []string) {
	switch v := in.(type) {
	case int:
		out = append(out, strconv.FormatInt(int64(v), 10))
	case []int:
		for _, vv := range v {
			out = append(out, strconv.FormatInt(int64(vv), 10))
		}
	case int8:
		out = append(out, strconv.FormatInt(int64(v), 10))
	case []int8:
		for _, vv := range v {
			out = append(out, strconv.FormatInt(int64(vv), 10))
		}
	case int16:
		out = append(out, strconv.FormatInt(int64(v), 10))
	case []int16:
		for _, vv := range v {
			out = append(out, strconv.FormatInt(int64(vv), 10))
		}
	case int32:
		out = append(out, strconv.FormatInt(int64(v), 10))
	case []int32:
		for _, vv := range v {
			out = append(out, strconv.FormatInt(int64(vv), 10))
		}
	case int64:
		out = append(out, strconv.FormatInt(v, 10))
	case []int64:
		for _, vv := range v {
			out = append(out, strconv.FormatInt(vv, 10))
		}
	case uint:
		out = append(out, strconv.FormatUint(uint64(v), 10))
	case []uint:
		for _, vv := range v {
			out = append(out, strconv.FormatUint(uint64(vv), 10))
		}
	case uint8:
		out = append(out, strconv.FormatUint(uint64(v), 10))
	case []uint8:
		for _, vv := range v {
			out = append(out, strconv.FormatUint(uint64(vv), 10))
		}
	case uint16:
		out = append(out, strconv.FormatUint(uint64(v), 10))
	case []uint16:
		for _, vv := range v {
			out = append(out, strconv.FormatUint(uint64(vv), 10))
		}
	case uint32:
		out = append(out, strconv.FormatUint(uint64(v), 10))
	case []uint32:
		for _, vv := range v {
			out = append(out, strconv.FormatUint(uint64(vv), 10))
		}
	case uint64:
		out = append(out, strconv.FormatUint(v, 10))
	case []uint64:
		for _, vv := range v {
			out = append(out, strconv.FormatUint(vv, 10))
		}
	case float32:
		out = append(out, strconv.FormatFloat(float64(v), 'f', -1, 32))
	case []float32:
		for _, vv := range v {
			out = append(out, strconv.FormatFloat(float64(vv), 'f', -1, 32))
		}
	case float64:
		out = append(out, strconv.FormatFloat(v, 'f', -1, 64))
	case []float64:
		for _, vv := range v {
			out = append(out, strconv.FormatFloat(vv, 'f', -1, 64))
		}
	case bool:
		out = append(out, strconv.FormatBool(v))
	case []bool:
		for _, vv := range v {
			out = append(out, strconv.FormatBool(vv))
		}
	case string:
		out = append(out, v)
	case []string:
		out = v
	case time.Time:
		out = append(out, v.Format(time.RFC3339))
	case []time.Time:
		for _, vv := range v {
			out = append(out, vv.Format(time.RFC3339))
		}
	case time.Duration:
		out = append(out, v.String())
	case []time.Duration:
		for _, vv := range v {
			out = append(out, vv.String())
		}
	default:
		out = append(out, fmt.Sprintf("%v", in))
	}
	return
}

// StructParams is a built-in implementation of ParamsCodec. It is mainly used
// to encode and decode a struct. The encoding and decoding of each field can
// be customized by setting Fields.
type StructParams struct {
	Fields map[string]ParamsCodec

	camelCase bool
}

func (p StructParams) CamelCase() StructParams {
	p.camelCase = true
	return p
}

// Decode decodes a string map to a struct (or a *struct).
func (p StructParams) Decode(in map[string][]string, out interface{}) error {
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

		kokField := &parser.StructField{
			Name:      field.Name,
			CamelCase: p.camelCase,
			Type:      field.Type.Name(),
			Tag:       field.Tag,
		}
		if err := kokField.Parse(); err != nil {
			return err
		}

		if kokField.Omitted {
			// Omit this field.
			continue
		}

		// Build values needed by this field for decoding.
		values := make(map[string][]string)
		for _, param := range kokField.Params {
			key := param.UniqueKey()
			values[key] = in[key]
		}

		fieldValuePtr := reflect.New(fieldValue.Type())

		codec := p.fieldCodec(field.Name)
		if err := codec.Decode(values, fieldValuePtr.Interface()); err != nil {
			return err
		}
		fieldValue.Set(fieldValuePtr.Elem())
	}

	return nil
}

func (p StructParams) fieldCodec(name string) ParamsCodec {
	if c, ok := p.Fields[name]; ok {
		return c
	}
	return defaultBasicParams
}

// Encode encodes a struct (or a *struct) to a string map.
func (p StructParams) Encode(in interface{}) (out map[string][]string) {
	inValue := reflect.ValueOf(in)
	switch k := inValue.Kind(); {
	case k == reflect.Ptr && inValue.Elem().Kind() == reflect.Struct:
		// Convert inValue from *struct to struct implicitly.
		inValue = inValue.Elem()
	case k == reflect.Struct:
	default:
		panic(ErrUnsupportedType)
	}

	outMap := make(map[string][]string)

	structType := inValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := inValue.Field(i)

		kokField := &parser.StructField{
			Name:      field.Name,
			CamelCase: p.camelCase,
			Type:      field.Type.Name(),
			Tag:       field.Tag,
		}
		if err := kokField.Parse(); err != nil {
			panic(err)
		}

		if kokField.Omitted {
			// Omit this field.
			continue
		}

		codec := p.fieldCodec(field.Name)
		for k, v := range codec.Encode(fieldValue.Interface()) {
			if k == specialKey {
				if len(kokField.Params) != 1 {
					panic(fmt.Errorf("special key %q is reserved for use in codecs wrapped by ToParamsCodec()", specialKey))
				}
				// Use the key of the parameter.
				k = kokField.Params[0].UniqueKey()
			}

			if k == "" {
				panic(fmt.Errorf("empty key returned from %s's codec", field.Name))
			}
			if _, ok := outMap[k]; ok {
				panic(fmt.Errorf("duplicate key %q returned from %s's codec", k, field.Name))
			}

			outMap[k] = v
		}
	}

	return outMap
}

// paramCodecAdapter turns the inner codec from a ParamCodec to a ParamsCodec.
type paramCodecAdapter struct {
	codec ParamCodec
}

// ToParamsCodec creates a ParamsCodec from a ParamCodec. It is mainly used
// along with StructParams.
func ToParamsCodec(codec ParamCodec) ParamsCodec {
	return &paramCodecAdapter{codec: codec}
}

func (p *paramCodecAdapter) Decode(in map[string][]string, out interface{}) error {
	if len(in) == 0 {
		return nil
	}

	if len(in) > 1 {
		return fmt.Errorf("only unable to decode multiple ")
	}

	var one []string
	for _, v := range in {
		one = v
		break
	}

	return p.codec.Decode(one, out)
}

func (p *paramCodecAdapter) Encode(in interface{}) map[string][]string {
	out := p.codec.Encode(in)
	return map[string][]string{specialKey: out}
}
