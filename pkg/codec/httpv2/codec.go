package codec

import (
	"io"
	"net/http"
)

type Codec interface {
	// Encoders and decoders used at the server side.
	DecodeRequestParam(name, value string, out interface{}) error
	DecodeRequestParams(name string, values map[string]string, out interface{}) error
	DecodeRequestBody(body io.ReadCloser, out interface{}) error
	EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error
	EncodeFailureResponse(w http.ResponseWriter, err error) error

	// Encoders and decoders used at the client side.
	EncodeRequestParam(name string, value interface{}) string
	EncodeRequestParams(name string, value interface{}) map[string]string
	EncodeRequestBody(body interface{}) (io.Reader, map[string]string, error)
	DecodeSuccessResponse(body io.ReadCloser, out interface{}) error
	DecodeFailureResponse(body io.ReadCloser, out *error) error
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}
