package oasv2

import (
	"net/http"
	"net/http/httptest"

	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"
)

type LegacyResponseSchema struct {
	Codecs          httpcodec.Codecs
	GetSuccessFunc  GetSuccessFunc
	GetFailuresFunc GetFailuresFunc
}

func (rs *LegacyResponseSchema) codecs() httpcodec.Codecs {
	if rs.Codecs == nil {
		rs.Codecs = httpcodec.CodecMap{}
	}
	return rs.Codecs
}

func (rs *LegacyResponseSchema) getSuccessFunc() GetSuccessFunc {
	if rs.GetSuccessFunc == nil {
		rs.GetSuccessFunc = func(body interface{}) interface{} {
			return body
		}
	}
	return rs.GetSuccessFunc
}

func (rs *LegacyResponseSchema) SuccessResponse(name string, statusCode int, body interface{}) Response {
	codec := rs.codecs().EncodeDecoder(name)

	w := httptest.NewRecorder()
	_ = codec.EncodeSuccessResponse(w, statusCode, body)
	contentType := w.Result().Header.Get("Content-Type")

	if statusCode == http.StatusNoContent || isMediaFile(contentType) {
		body = nil
	} else {
		body = rs.getSuccessFunc()(body)
	}

	return Response{
		StatusCode:  statusCode,
		ContentType: contentType,
		Body:        body,
	}
}

func (rs *LegacyResponseSchema) FailureResponses(name string) (resps []Response) {
	if rs.GetFailuresFunc == nil {
		return
	}

	codec := rs.codecs().EncodeDecoder(name)

	for err, body := range rs.GetFailuresFunc(name) {
		w := httptest.NewRecorder()
		_ = codec.EncodeFailureResponse(w, err)

		contentType := w.Result().Header.Get("Content-Type")
		if body == nil {
			body = decodePerContentType(contentType, w.Result().Body)
		}

		resps = append(resps, Response{
			StatusCode:  w.Result().StatusCode,
			ContentType: contentType,
			Body:        body,
		})
	}

	return
}
