package profilesvc

import (
	"net/http"

	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"
	"github.com/RussellLuo/kok/pkg/oasv2"
)

type Codec struct {
	httpcodec.JSONCodec
}

func (c Codec) EncodeFailureResponse(w http.ResponseWriter, err error) error {
	return c.JSONCodec.EncodeSuccessResponse(w, codeFrom(err), toBody(err))
}

func toBody(err error) interface{} {
	return map[string]string{
		"error": err.Error(),
	}
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func NewCodecs() httpcodec.Codecs {
	return httpcodec.CodecMap{
		Default: Codec{},
	}
}

func GetFailures(name string) map[error]interface{} {
	toFailures := func(errs ...error) map[error]interface{} {
		m := make(map[error]interface{})
		for _, err := range errs {
			m[err] = toBody(err)
		}
		return m
	}

	switch name {
	case "PostProfile":
		return toFailures(ErrAlreadyExists)
	case "GetProfile", "DeleteProfile", "GetAddresses", "GetAddress", "DeleteAddress":
		return toFailures(ErrNotFound)
	case "PutProfile":
		return toFailures(ErrInconsistentIDs)
	case "PatchProfile":
		return toFailures(ErrInconsistentIDs, ErrNotFound)
	case "PostAddress":
		return toFailures(ErrAlreadyExists, ErrNotFound)
	default:
		return nil
	}
}

func NewSchema() oasv2.Schema {
	return &oasv2.ResponseSchema{
		Codecs:          NewCodecs(),
		GetFailuresFunc: GetFailures,
	}
}
