package codec

import (
	"fmt"
	"strconv"
	"time"
)

type ParamCodec struct {
	OnDecode func(value string) (interface{}, error)
	OnEncode func(value interface{}) string
}

func (pc ParamCodec) Decode(name, value string, out interface{}) error {
	switch t := out.(type) {
	case *int:
		return pc.decodeInt(name, value, out.(*int))
	case *int8:
		return pc.decodeInt8(name, value, out.(*int8))
	case *int16:
		return pc.decodeInt16(name, value, out.(*int16))
	case *int32:
		return pc.decodeInt32(name, value, out.(*int32))
	case *int64:
		return pc.decodeInt64(name, value, out.(*int64))
	case *uint:
		return pc.decodeUint(name, value, out.(*uint))
	case *uint8:
		return pc.decodeUint8(name, value, out.(*uint8))
	case *uint16:
		return pc.decodeUint16(name, value, out.(*uint16))
	case *uint32:
		return pc.decodeUint32(name, value, out.(*uint32))
	case *uint64:
		return pc.decodeUint64(name, value, out.(*uint64))
	case *bool:
		return pc.decodeBool(name, value, out.(*bool))
	case *string:
		return pc.decodeString(name, value, out.(*string))
	case *time.Time:
		return pc.decodeTime(name, value, out.(*time.Time))
	default:
		// Panic since this is a programming error.
		panic(fmt.Errorf("unsupported out type: %T", t))
	}

	return nil
}

func (pc ParamCodec) Encode(name string, value interface{}) string {
	if pc.OnEncode != nil {
		return pc.OnEncode(value)
	}

	switch value.(type) {
	case int:
		v := value.(int)
		return strconv.FormatInt(int64(v), 10)
	case int8:
		v := value.(int8)
		return strconv.FormatInt(int64(v), 10)
	case int16:
		v := value.(int16)
		return strconv.FormatInt(int64(v), 10)
	case int32:
		v := value.(int32)
		return strconv.FormatInt(int64(v), 10)
	case int64:
		v := value.(int64)
		return strconv.FormatInt(v, 10)
	case uint:
		v := value.(uint)
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		v := value.(uint8)
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		v := value.(uint16)
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		v := value.(uint32)
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		v := value.(uint64)
		return strconv.FormatUint(v, 10)
	case bool:
		v := value.(bool)
		return strconv.FormatBool(v)
	case string:
		return value.(string)
	case time.Time:
		v := value.(time.Time)
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (pc ParamCodec) decodeInt(name, value string, out *int) error {
	if pc.OnDecode == nil {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*out = v
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(int)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want int)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeInt8(name, value string, out *int8) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = int8(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(int8)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want int8)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeInt16(name, value string, out *int16) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = int16(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(int16)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want int16)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeInt32(name, value string, out *int32) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = int32(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(int32)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want int32)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeInt64(name, value string, out *int64) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out = v
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(int64)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want int64)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeUint(name, value string, out *uint) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = uint(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(uint)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want uint)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeUint8(name, value string, out *uint8) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = uint8(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(uint8)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want uint8)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeUint16(name, value string, out *uint16) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = uint16(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(uint16)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want uint16)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeUint32(name, value string, out *uint32) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = uint32(v)
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(uint32)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want uint32)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeUint64(name, value string, out *uint64) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out = v
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(uint64)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want uint64)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeBool(name, value string, out *bool) error {
	if pc.OnDecode == nil {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*out = v
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(bool)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want bool)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeString(name, value string, out *string) error {
	if pc.OnDecode == nil {
		*out = value
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(string)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want string)", name, result))
	}

	*out = v
	return nil
}

func (pc ParamCodec) decodeTime(name, value string, out *time.Time) error {
	if pc.OnDecode == nil {
		v, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		*out = v
		return nil
	}

	result, err := pc.OnDecode(value)
	if err != nil {
		return err
	}

	v, ok := result.(time.Time)
	if !ok {
		// Panic since this is a programming error.
		panic(fmt.Errorf("decoder of %q returns %v (want time.Time)", name, result))
	}

	*out = v
	return nil
}
