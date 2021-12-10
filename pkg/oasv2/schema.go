package oasv2

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/RussellLuo/kun/pkg/codec/httpcodec"
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
		rs.Codecs = httpcodec.NewDefaultCodecs(nil)
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

func (rs *ResponseSchema) FailureResponses(name string) (resps []Response) {
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

func decodePerContentType(contentType string, body io.ReadCloser) (out map[string]interface{}) {
	switch {
	case strings.HasPrefix(contentType, "application/json"):
		_ = httpcodec.JSON{}.DecodeSuccessResponse(body, &out)
	}
	return
}

func Errors(errs ...error) map[error]interface{} {
	m := make(map[error]interface{})
	for _, err := range errs {
		m[err] = nil
	}
	return m
}
