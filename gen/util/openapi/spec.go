package openapi

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/RussellLuo/kok/gen/util/misc"
)

const (
	InPath   = "path"
	InQuery  = "query"
	InHeader = "header"
	InCookie = "cookie"
	InBody   = "body"

	MediaTypeJSON = "application/json; charset=utf-8"
)

type Specification struct {
	Operations []*Operation
}

func Spec() *Specification {
	return &Specification{}
}

func (s *Specification) Path(pattern string, operations ...*Operation) *Specification {
	for _, o := range operations {
		o.Pattern = pattern
	}
	s.Operations = append(s.Operations, operations...)
	return s
}

type Param struct {
	Name     string
	In       string
	Alias    string
	Type     string
	Required bool

	// The name of the decoder function, which is used to convert the value
	// from string to another type (e.g. integer, boolean or struct).
	Decoder string

	Sub []*Param
}

func (p *Param) SetName(name string) {
	p.Name = name

	// Set alias if it's not set.
	if p.Alias == "" {
		p.Alias = name
	}
}

// Set sets properties according to the values hold by o.
func (p *Param) Set(o *Param) {
	if len(p.Sub) != 0 {
		panic(fmt.Errorf("parent param %q can not be used alone", p.Name))
	}

	/*if !isPrimitiveType(p.Type) && o.In != InBody {
		panic(fmt.Errorf("non-primitive param %q must be in `body`", p.Name))
	}*/

	p.In = o.In
	p.Alias = o.Alias
	p.Required = o.Required
	p.Decoder = o.Decoder
}

// Add adds o as a sub parameter of the current parameter.
func (p *Param) Add(o *Param) {
	if isPrimitiveType(p.Type) {
		panic(fmt.Errorf("primitive param %q can not has sub parameters", p.Name))
	}

	p.Sub = append(p.Sub, o)

	// Clear the properties that are meaningless for the parent parameter.
	p.In = ""
	p.Alias = ""
	p.Required = false
}

type Request struct {
	MediaType string
	Params    []*Param
}

type Response struct {
	StatusCode int
	MediaType  string
	Schema     interface{}
}

type Options struct {
	/*RequestDecoder struct {
		// The name of the decoder function, which is used to extract the parameters
		// in body by decoding the request body.
		Body string
	}*/

	ResponseEncoder struct {
		// The name of the encoder function for the successful response.
		Success string
		// The name of the encoder function for the error response.
		Failure string
	}
}

type Operation struct {
	Name             string
	Method           string
	Pattern          string
	Request          Request
	SuccessResponse  *Response
	FailureResponses []*Response
	Options          Options
}

func GET() *Operation {
	return &Operation{Method: "GET"}
}

func POST() *Operation {
	return &Operation{Method: "POST"}
}

func PUT() *Operation {
	return &Operation{Method: "PUT"}
}

func PATCH() *Operation {
	return &Operation{Method: "PATCH"}
}

func DELETE() *Operation {
	return &Operation{Method: "DELETE"}
}

func OPTIONS() *Operation {
	return &Operation{Method: "OPTIONS"}
}

func HEAD() *Operation {
	return &Operation{Method: "HEAD"}
}

func (o *Operation) Req(mediaType string, schema interface{}) *Operation {
	o.Request.MediaType = mediaType

	t := reflect.TypeOf(schema)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		p := o.buildParam(field.Tag.Get("kok"), field.Name, field.Type.String())
		o.addParam(p)
	}

	return o
}

func (o *Operation) buildParam(text, name, typ string) *Param {
	p := &Param{Type: typ}

	for _, part := range strings.Split(text, ",") {
		if !strings.Contains(part, ":") {
			panic(fmt.Errorf("invalid tag part: %s", part))
		}

		split := strings.SplitN(part, ":", 2)
		key, value := split[0], split[1]

		switch key {
		case "type":
			p.Type = value
		case "in":
			p.In = value
			if value == InPath {
				p.Required = true
			}
		case "required":
			p.Required = value == "true"
		case "name":
			p.SetName(value)
		case "alias":
			p.Alias = value
		case "decoder":
			p.Decoder = value
		default:
			panic(fmt.Errorf("invalid tag part: %s", part))
		}
	}

	if p.In == "" {
		p.In = InBody
	}
	if p.Name == "" && name != "" {
		p.SetName(misc.LowerFirst(name))
	}

	if strings.Contains(p.Name, ".") && p.In == InBody {
		panic(fmt.Errorf("sub param %q must be in `path`, `query` or `header`", p.Name))
	}

	return p
}

func (o *Operation) addParam(p *Param) *Operation {
	for _, param := range o.Request.Params {
		if p.Name == param.Name {
			panic(errors.New("duplicate parameter name " + p.Name))
		}
	}

	o.Request.Params = append(o.Request.Params, p)
	return o
}

func (o *Operation) Resp(statusCode int, mediaType string, schema interface{}) *Operation {
	if mediaType != MediaTypeJSON && !strings.HasPrefix(mediaType, "image/") {
		panic(errors.New(mediaType + " not supported"))
	}

	if statusCode >= http.StatusContinue && statusCode < http.StatusBadRequest {
		o.SuccessResponse = &Response{
			StatusCode: statusCode,
			MediaType:  mediaType,
			Schema:     schema,
		}
	} else {
		o.FailureResponses = append(o.FailureResponses, &Response{
			StatusCode: statusCode,
			MediaType:  mediaType,
			Schema:     schema,
		})
	}

	return o
}

func (o *Operation) Alias(name string) *Operation {
	o.Name = name
	return o
}
