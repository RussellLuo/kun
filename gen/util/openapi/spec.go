package openapi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	InPath   = "path"
	InQuery  = "query"
	InHeader = "header"
	InCookie = "cookie"
	InBody   = "body"

	MediaTypeJSON = "application/json"
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
	In       string
	Name     string
	Alias    string
	Type     string
	Required bool
}

func (p *Param) SetName(name string) {
	p.Name = name

	// Set alias if it's not set.
	if p.Alias == "" {
		p.Alias = name
	}
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
	ErrorEncoder string
}

type Operation struct {
	Name      string
	Method    string
	Pattern   string
	Request   Request
	Responses []Response
	Options   Options
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
		default:
			panic(fmt.Errorf("invalid tag part: %s", part))
		}
	}

	if p.In == "" {
		p.In = InBody
	}
	if p.Name == "" && name != "" {
		p.SetName(strings.ToLower(string(name[0])) + name[1:])
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
	if mediaType != MediaTypeJSON {
		panic(errors.New(mediaType + " not supported"))
	}

	o.Responses = append(o.Responses, Response{
		StatusCode: statusCode,
		MediaType:  mediaType,
		Schema:     schema,
	})
	return o
}

func (o *Operation) Alias(name string) *Operation {
	o.Name = name
	return o
}
