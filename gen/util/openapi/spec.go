package openapi

import (
	"errors"
	"fmt"
	"go/types"
	"net/http"
	"reflect"
	"strings"

	"github.com/RussellLuo/kun/pkg/caseconv"
)

const (
	InPath   = "path"
	InQuery  = "query"
	InHeader = "header"
	InCookie = "cookie"
	InBody   = "body"

	// This flag indicates that the parameter is located in *http.Request.
	InRequest = "request"

	MediaTypeJSON = "application/json; charset=utf-8"
)

type Specification struct {
	Metadata   *Metadata
	Operations []*Operation
}

type Metadata struct {
	// Relative path to the OAS documentation.
	DocsPath string

	Title       string
	Version     string
	Description string
	BasePath    string

	// Default tags for operations those have no tags.
	DefaultTags []string
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
	Name        string     // Method argument name
	Type        string     // Method argument type
	RawType     types.Type // The raw Go type of the method argument
	In          string
	Alias       string // Request parameter name
	AliasType   string // Request parameter type
	Required    bool
	Description string // OAS description

	IsBlank bool // Whether this parameter is a blank identifier.

	inUse bool // Indicates this parameter already has a corresponding @kok(param).
}

func (p *Param) SetName(name string, snakeCase bool) {
	p.Name = name

	// Set alias if it's not set.
	if p.Alias == "" {
		if snakeCase {
			p.Alias = caseconv.ToSnakeCase(name)
		} else {
			p.Alias = caseconv.ToLowerCamelCase(name)
		}
	}
}

func (p *Param) SetAlias(alias string) {
	p.Alias = alias
}

func (p *Param) SetDescription(descr string) {
	p.Description = descr
}

// Set sets properties according to the values hold by o.
func (p *Param) Set(o *Param) {
	/*if !isPrimitiveType(p.Type) && o.In != InBody {
		panic(fmt.Errorf("non-primitive param %q must be in `body`", p.Name))
	}*/

	p.In = o.In
	p.Alias = o.Alias
	if o.AliasType != "" {
		p.AliasType = o.AliasType
	}
	p.Required = o.Required

	p.inUse = true
}

// Set sets properties according to the given annotation.
func (p *Param) SetByAnnotation(a *annotation) {
	if a.In != "" {
		p.In = a.In
	}
	if a.Name != "" {
		p.Alias = a.Name
	}
	p.Required = a.Required
	p.AliasType = a.Type
	p.Description = a.Description

	p.inUse = true
}

type Request struct {
	MediaType string

	// The name of the request field whose value is mapped to the HTTP request body.
	// Otherwise, all fields not located in path/query/header will be mapped to the HTTP body
	BodyField string

	Params []*Param
}

type Response struct {
	StatusCode int
	MediaType  string
	Schema     interface{}

	// The name of the response field whose value is mapped to the HTTP response body.
	// When omitted, the entire response struct will be used as the HTTP response body.
	BodyField string
}

type Operation struct {
	Name             string
	Method           string
	Pattern          string
	Request          Request
	SuccessResponse  *Response
	FailureResponses []*Response
	Description      string
	Tags             []string
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

// DEPRECATED?
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
			p.SetName(value, false)
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
		p.SetName(caseconv.LowerFirst(name), false)
	}

	if p.In == InRequest && p.Alias != "RemoteAddr" {
		panic(fmt.Errorf("param %q tries to extract value from `request.%s`, but only `request.RemoteAddr` is available", p.Name, p.Alias))
	}

	if strings.Contains(p.Name, ".") && p.In == InBody {
		panic(fmt.Errorf("sub param %q must be in `path`, `query`, `header` or `request`", p.Name))
	}

	return p
}

func (o *Operation) addParam(p *Param) *Operation {
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
