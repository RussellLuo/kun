package httpcodec

import (
	"io"
	"net/http"
)

// Codec is a series of codecs (encoders and decoders) for HTTP requests and responses.
type Codec interface {
	// Encoders and decoders used at the server side.
	DecodeRequestParam(name string, values []string, out interface{}) error
	DecodeRequestParams(name string, values map[string][]string, out interface{}) error
	DecodeRequestBody(r *http.Request, out interface{}) error
	EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error
	EncodeFailureResponse(w http.ResponseWriter, err error) error

	// Encoders and decoders used at the client side.
	EncodeRequestParam(name string, value interface{}) []string
	EncodeRequestParams(name string, value interface{}) map[string][]string
	EncodeRequestBody(body interface{}) (io.Reader, map[string]string, error)
	DecodeSuccessResponse(body io.ReadCloser, out interface{}) error
	DecodeFailureResponse(body io.ReadCloser, out *error) error
}

type Codecs interface {
	EncodeDecoder(name string) Codec
}

// ParamCodec is a codec (encoder and decoder) for a single request parameter.
type ParamCodec interface {
	Decode(in []string, out interface{}) error
	Encode(in interface{}) (out []string)
}

// ParamsCodec is a codec (encoder and decoder) for a group of request parameters.
type ParamsCodec interface {
	Decode(in map[string][]string, out interface{}) error
	Encode(in interface{}) (out map[string][]string)
}
