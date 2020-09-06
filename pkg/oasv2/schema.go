package oasv2

import (
	"net/http/httptest"

	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"
)

type Response struct {
	StatusCode  int
	ContentType string
	Body        interface{}
}

type Schema interface {
	SuccessResponse(name string, statusCode int, body interface{}) Response
	FailureResponses(name string) []Response
}

type GetSuccessFunc func(body interface{}) interface{}
type GetFailuresFunc func(name string) map[error]interface{}

type ResponseSchema struct {
	Codecs          httpcodec.Codecs
	GetSuccessFunc  GetSuccessFunc
	GetFailuresFunc GetFailuresFunc
}

func (rs *ResponseSchema) codecs() httpcodec.Codecs {
	if rs.Codecs == nil {
		rs.Codecs = httpcodec.CodecMap{}
	}
	return rs.Codecs
}

func (rs *ResponseSchema) getSuccessFunc() GetSuccessFunc {
	if rs.GetSuccessFunc == nil {
		rs.GetSuccessFunc = func(body interface{}) interface{} {
			return body
		}
	}
	return rs.GetSuccessFunc
}

func (rs *ResponseSchema) SuccessResponse(name string, statusCode int, body interface{}) Response {
	codec := rs.codecs().EncodeDecoder(name)

	w := httptest.NewRecorder()
	_ = codec.EncodeSuccessResponse(w, statusCode, body)
	contentType := w.Result().Header.Get("Content-Type")

	if !isMediaFile(contentType) {
		body = rs.getSuccessFunc()(body)
	} else {
		body = nil
	}

	return Response{
		StatusCode:  statusCode,
		ContentType: contentType,
		Body:        body,
	}
}

func (rs *ResponseSchema) FailureResponses(name string) (resps []Response) {
	if rs.GetFailuresFunc == nil {
		return
	}

	codec := rs.codecs().EncodeDecoder(name)

	for err, body := range rs.GetFailuresFunc(name) {
		w := httptest.NewRecorder()
		_ = codec.EncodeFailureResponse(w, err)
		resps = append(resps, Response{
			StatusCode:  w.Result().StatusCode,
			ContentType: w.Result().Header.Get("Content-Type"),
			Body:        body,
		})
	}

	return
}
