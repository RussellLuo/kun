package codec

import (
	"fmt"
	"strconv"
	"time"
)

type Param struct{}

func (p Param) Decode(name string, values []string, out interface{}) error {
	if len(values) == 0 {
		// For legacy generated code, values may be an empty slice when the parameter is optional.
		return nil
	}

	switch v := out.(type) {
	case *int:
		return p.decodeInt(name, values[0], v)
	case *[]int:
		return p.decodeIntSlice(name, values, v)
	case *int8:
		return p.decodeInt8(name, values[0], v)
	case *[]int8:
		return p.decodeInt8Slice(name, values, v)
	case *int16:
		return p.decodeInt16(name, values[0], v)
	case *[]int16:
		return p.decodeInt16Slice(name, values, v)
	case *int32:
		return p.decodeInt32(name, values[0], v)
	case *[]int32:
		return p.decodeInt32Slice(name, values, v)
	case *int64:
		return p.decodeInt64(name, values[0], v)
	case *[]int64:
		return p.decodeInt64Slice(name, values, v)
	case *uint:
		return p.decodeUint(name, values[0], v)
	case *[]uint:
		return p.decodeUintSlice(name, values, v)
	case *uint8:
		return p.decodeUint8(name, values[0], v)
	case *[]uint8:
		return p.decodeUint8Slice(name, values, v)
	case *uint16:
		return p.decodeUint16(name, values[0], v)
	case *[]uint16:
		return p.decodeUint16Slice(name, values, v)
	case *uint32:
		return p.decodeUint32(name, values[0], v)
	case *[]uint32:
		return p.decodeUint32Slice(name, values, v)
	case *uint64:
		return p.decodeUint64(name, values[0], v)
	case *[]uint64:
		return p.decodeUint64Slice(name, values, v)
	case *bool:
		return p.decodeBool(name, values[0], v)
	case *[]bool:
		return p.decodeBoolSlice(name, values, v)
	case *string:
		return p.decodeString(name, values[0], v)
	case *[]string:
		return p.decodeStringSlice(name, values, v)
	case *time.Time:
		return p.decodeTime(name, values[0], v)
	case *[]time.Time:
		return p.decodeTimeSlice(name, values, v)
	case *time.Duration:
		return p.decodeDuration(name, values[0], v)
	case *[]time.Duration:
		return p.decodeDurationSlice(name, values, v)
	default:
		// Panic since this is a programming error.
		panic(fmt.Errorf("unsupported out type: %T", v))
	}
}

func (p Param) Encode(name string, in interface{}) (values []string) {
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

func (p Param) decodeInt(name, value string, out *int) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeIntSlice(name string, values []string, out *[]int) error {
	for _, value := range values {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return nil
}

func (p Param) decodeInt8(name, value string, out *int8) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*out = int8(v)
	return nil
}

func (p Param) decodeInt8Slice(name string, values []string, out *[]int8) error {
	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, int8(v))
	}
	return nil
}

func (p Param) decodeInt16(name, value string, out *int16) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*out = int16(v)
	return nil
}

func (p Param) decodeInt16Slice(name string, values []string, out *[]int16) error {
	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, int16(v))
	}
	return nil
}

func (p Param) decodeInt32(name, value string, out *int32) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*out = int32(v)
	return nil
}

func (p Param) decodeInt32Slice(name string, values []string, out *[]int32) error {
	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, int32(v))
	}
	return nil
}

func (p Param) decodeInt64(name, value string, out *int64) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeInt64Slice(name string, values []string, out *[]int64) error {
	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, int64(v))
	}
	return nil
}

func (p Param) decodeUint(name, value string, out *uint) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*out = uint(v)
	return nil
}

func (p Param) decodeUintSlice(name string, values []string, out *[]uint) error {
	for _, value := range values {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, uint(v))
	}
	return nil
}

func (p Param) decodeUint8(name, value string, out *uint8) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*out = uint8(v)
	return nil
}

func (p Param) decodeUint8Slice(name string, values []string, out *[]uint8) error {
	for _, value := range values {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, uint8(v))
	}
	return nil
}

func (p Param) decodeUint16(name, value string, out *uint16) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*out = uint16(v)
	return nil
}

func (p Param) decodeUint16Slice(name string, values []string, out *[]uint16) error {
	for _, value := range values {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, uint16(v))
	}
	return nil
}

func (p Param) decodeUint32(name, value string, out *uint32) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*out = uint32(v)
	return nil
}

func (p Param) decodeUint32Slice(name string, values []string, out *[]uint32) error {
	for _, value := range values {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, uint32(v))
	}
	return nil
}

func (p Param) decodeUint64(name, value string, out *uint64) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeUint64Slice(name string, values []string, out *[]uint64) error {
	for _, value := range values {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return nil
}

func (p Param) decodeBool(name, value string, out *bool) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeBoolSlice(name string, values []string, out *[]bool) error {
	for _, value := range values {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return nil
}

func (p Param) decodeString(name, value string, out *string) error {
	*out = value
	return nil
}

func (p Param) decodeStringSlice(name string, values []string, out *[]string) error {
	*out = append(*out, values...)
	return nil
}

func (p Param) decodeTime(name, value string, out *time.Time) error {
	v, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeTimeSlice(name string, values []string, out *[]time.Time) error {
	for _, value := range values {
		v, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return nil
}

func (p Param) decodeDuration(name, value string, out *time.Duration) error {
	v, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*out = v
	return nil
}

func (p Param) decodeDurationSlice(name string, values []string, out *[]time.Duration) error {
	for _, value := range values {
		v, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*out = append(*out, v)
	}
	return nil
}
