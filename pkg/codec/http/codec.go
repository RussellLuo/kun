package codec

import (
	"io"
	"net/http"
)

type Codec interface {
	DecodeRequestParam(name, value string, out interface{}) error
	DecodeRequestBody(body io.ReadCloser, out interface{}) error
	EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error
	EncodeFailureResponse(w http.ResponseWriter, err error) error
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}
