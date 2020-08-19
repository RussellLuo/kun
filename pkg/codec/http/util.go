package codec

import (
	"fmt"
	"strconv"
	"time"
)

func DecodeStringPerOutType(value string, out interface{}) error {
	switch t := out.(type) {
	case *int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*out.(*int) = v
	case *int8:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*int8) = int8(v)
	case *int16:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*int16) = int16(v)
	case *int32:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*int32) = int32(v)
	case *int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*int64) = v
	case *uint:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*uint) = uint(v)
	case *uint8:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*uint8) = uint8(v)
	case *uint16:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*uint16) = uint16(v)
	case *uint32:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*uint32) = uint32(v)
	case *uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		*out.(*uint64) = v
	case *bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*out.(*bool) = v
	case *string:
		*out.(*string) = value
	case *time.Time:
		v, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		*out.(*time.Time) = v
	default:
		// Panic since this is a programming error.
		panic(fmt.Errorf("unsupported type %v", t))
	}

	return nil
}
