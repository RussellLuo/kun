package codec

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type Bodier interface {
	Body() interface{}
}

func MakeResponseEncoder(codec Codec, statusCode int) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			return f.Failed()
		}

		if statusCode == http.StatusNoContent {
			// Respond with no content.
			w.WriteHeader(statusCode)
			return nil
		}

		body, ok := response.(Bodier)
		if ok {
			return codec.EncodeSuccessResponse(w, statusCode, body.Body())
		}
		return codec.EncodeSuccessResponse(w, statusCode, response)
	}
}

func MakeErrorEncoder(codec Codec) kithttp.ErrorEncoder {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		_ = codec.EncodeFailureResponse(w, err)
	}
}
