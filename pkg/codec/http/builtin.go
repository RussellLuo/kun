package codec

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RussellLuo/kok/pkg/werror/googlecode"
)

type CodecMap struct {
	Codecs  map[string]Codec
	Default Codec
}

func (cm CodecMap) EncodeDecoder(name string) Codec {
	if c, ok := cm.Codecs[name]; ok {
		return c
	}

	if cm.Default != nil {
		return cm.Default
	}
	return NewJSONCodec(nil) // defaults to JSONCodec
}

type ParamDecoder func(value string) (interface{}, error)

type JSONCodec struct {
	paramDecoders map[string]ParamDecoder
}

func NewJSONCodec(paramDecoders map[string]ParamDecoder) JSONCodec {
	return JSONCodec{paramDecoders: paramDecoders}
}

func (jc JSONCodec) DecodeRequestParam(name, value string, out interface{}) error {
	decoder, ok := jc.paramDecoders[name]
	if !ok {
		return DecodeStringPerOutType(value, out)
	}

	result, err := decoder(value)
	if err != nil {
		return err
	}

	switch t := out.(type) {
	case *int:
		v, ok := result.(int)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want int)", name, result))
		}
		*out.(*int) = v
	case *int8:
		v, ok := result.(int8)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want int8)", name, result))
		}
		*out.(*int8) = v
	case *int16:
		v, ok := result.(int16)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want int16)", name, result))
		}
		*out.(*int16) = v
	case *int32:
		v, ok := result.(int32)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want int32)", name, result))
		}
		*out.(*int32) = v
	case *int64:
		v, ok := result.(int64)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want int64)", name, result))
		}
		*out.(*int64) = v
	case *uint:
		v, ok := result.(uint)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want uint)", name, result))
		}
		*out.(*uint) = v
	case *uint8:
		v, ok := result.(uint8)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want uint8)", name, result))
		}
		*out.(*uint8) = v
	case *uint16:
		v, ok := result.(uint16)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want uint16)", name, result))
		}
		*out.(*uint16) = v
	case *uint32:
		v, ok := result.(uint32)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want uint32)", name, result))
		}
		*out.(*uint32) = v
	case *uint64:
		v, ok := result.(uint64)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want uint64)", name, result))
		}
		*out.(*uint64) = v
	case *bool:
		v, ok := result.(bool)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want bool)", name, result))
		}
		*out.(*bool) = v
	case *string:
		v, ok := result.(string)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want string)", name, result))
		}
		*out.(*string) = v
	case *time.Time:
		v, ok := result.(time.Time)
		if !ok {
			// Panic since this is a programming error.
			panic(fmt.Errorf("decoder of %q returns %v (want time.Time)", name, result))
		}
		*out.(*time.Time) = v
	default:
		// Panic since this is a programming error.
		panic(fmt.Errorf("unsupported out type: %T", t))
	}

	return nil
}

func (jc JSONCodec) DecodeRequestBody(body io.ReadCloser, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}

func (jc JSONCodec) EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

func (jc JSONCodec) EncodeFailureResponse(w http.ResponseWriter, err error) error {
	statusCode := googlecode.HTTPStatusCode(err)
	code, message := googlecode.ToCodeMessage(err)
	return jc.EncodeSuccessResponse(w, statusCode, map[string]map[string]string{
		"error": {
			"code":    code,
			"message": message,
		},
	})
}
