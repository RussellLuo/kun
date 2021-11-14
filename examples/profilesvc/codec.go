package profilesvc

import (
	"net/http"

	"github.com/RussellLuo/kok/pkg/httpcodec"
	"github.com/RussellLuo/kok/pkg/oas2"
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

func NewCodecs() *httpcodec.DefaultCodecs {
	return httpcodec.NewDefaultCodecs(Codec{})
}

func GetFailures(name string) map[error]interface{} {
	switch name {
	case "PostProfile":
		return oas2.Errors(ErrAlreadyExists)
	case "GetProfile", "DeleteProfile", "GetAddresses", "GetAddress", "DeleteAddress":
		return oas2.Errors(ErrNotFound)
	case "PutProfile":
		return oas2.Errors(ErrInconsistentIDs)
	case "PatchProfile":
		return oas2.Errors(ErrInconsistentIDs, ErrNotFound)
	case "PostAddress":
		return oas2.Errors(ErrAlreadyExists, ErrNotFound)
	default:
		return nil
	}
}

func NewSchema() *oas2.ResponseSchema {
	return &oas2.ResponseSchema{
		Codecs:          NewCodecs(),
		GetFailuresFunc: GetFailures,
	}
}
