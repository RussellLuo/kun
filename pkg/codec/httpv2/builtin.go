package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RussellLuo/kok/pkg/werror"
	"github.com/RussellLuo/kok/pkg/werror/gcode"
)

type CodecMap struct {
	Codecs  map[string]Codec
	Default Codec
}

func (cm CodecMap) EncodeDecoder(name string) Codec {
	if c, ok := cm.Codecs[name]; ok {
		return c
	}

	if cm.Default != nil {
		return cm.Default
	}
	return NewJSONCodec(nil) // defaults to JSONCodec
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FailureResponse struct {
	Error Error `json:"error"`
}

type JSONCodec struct {
	paramCodecs       map[string]ParamCodec
	defaultParamCodec ParamCodec
}

func NewJSONCodec(paramCodecs map[string]ParamCodec) JSONCodec {
	return JSONCodec{
		paramCodecs:       paramCodecs,
		defaultParamCodec: ParamCodec{},
	}
}

func (jc JSONCodec) DecodeRequestParam(name, value string, out interface{}) error {
	pc, ok := jc.paramCodecs[name]
	if !ok {
		pc = jc.defaultParamCodec
	}

	if err := pc.Decode(name, value, out); err != nil {
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (jc JSONCodec) DecodeRequestParams(name string, values map[string]string, out interface{}) error {
	if err := DecodeMapToStruct(values, out); err != nil {
		if err == ErrUnsupportedType {
			panic(fmt.Errorf("DecodeRequestParams not implemented for %q (of type %T)", name, out))
		}
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (jc JSONCodec) DecodeRequestBody(body io.ReadCloser, out interface{}) error {
	if err := json.NewDecoder(body).Decode(out); err != nil {
		return werror.Wrap(gcode.ErrInvalidArgument, err)
	}
	return nil
}

func (jc JSONCodec) EncodeSuccessResponse(w http.ResponseWriter, statusCode int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

func (jc JSONCodec) EncodeFailureResponse(w http.ResponseWriter, err error) error {
	statusCode := gcode.HTTPStatusCode(err)
	code, message := gcode.ToCodeMessage(err)
	return jc.EncodeSuccessResponse(w, statusCode, FailureResponse{
		Error: Error{
			Code:    code,
			Message: message,
		},
	})
}

func (jc JSONCodec) EncodeRequestParam(name string, value interface{}) string {
	pc, ok := jc.paramCodecs[name]
	if !ok {
		pc = jc.defaultParamCodec
	}
	return pc.Encode(name, value)
}

func (jc JSONCodec) EncodeRequestParams(name string, value interface{}) map[string]string {
	out := make(map[string]string)
	if err := EncodeStructToMap(value, &out); err != nil {
		if err == ErrUnsupportedType {
			panic(fmt.Errorf("EncodeRequestParams not implemented for %q (of type %T)", name, out))
		}
		panic(err)
	}
	return nil
}

func (jc JSONCodec) EncodeRequestBody(body interface{}) (io.Reader, map[string]string, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
	return bytes.NewBuffer(data), headers, nil
}

func (jc JSONCodec) DecodeSuccessResponse(body io.ReadCloser, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}

func (jc JSONCodec) DecodeFailureResponse(body io.ReadCloser, out *error) error {
	var resp FailureResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return err
	}

	*out = gcode.FromCodeMessage(resp.Error.Code, resp.Error.Message)
	return nil
}
