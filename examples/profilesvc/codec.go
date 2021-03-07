package profilesvc

import (
	"net/http"

	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
	"github.com/RussellLuo/kok/pkg/oasv2"
)

type Codec struct {
	httpcodec.JSON
}

func (c Codec) EncodeFailureResponse(w http.ResponseWriter, err error) error {
	return c.JSON.EncodeSuccessResponse(w, codeFrom(err), map[string]string{
		"error": err.Error(),
	})
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

func NewCodecs() httpcodec.CodecMap {
	return httpcodec.CodecMap{
		Default: Codec{},
	}
}

func GetFailures(name string) map[error]interface{} {
	switch name {
	case "PostProfile":
		return oasv2.Errors(ErrAlreadyExists)
	case "GetProfile", "DeleteProfile", "GetAddresses", "GetAddress", "DeleteAddress":
		return oasv2.Errors(ErrNotFound)
	case "PutProfile":
		return oasv2.Errors(ErrInconsistentIDs)
	case "PatchProfile":
		return oasv2.Errors(ErrInconsistentIDs, ErrNotFound)
	case "PostAddress":
		return oasv2.Errors(ErrAlreadyExists, ErrNotFound)
	default:
		return nil
	}
}

func NewSchema() *oasv2.ResponseSchema {
	return &oasv2.ResponseSchema{
		Codecs:          NewCodecs(),
		GetFailuresFunc: GetFailures,
	}
}
