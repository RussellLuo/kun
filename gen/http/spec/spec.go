package spec

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/ifacetool"
)

type Location string

const (
	InPath   Location = "path"
	InQuery  Location = "query"
	InHeader Location = "header"
	InCookie Location = "cookie"
	InBody   Location = "body"

	// InRequest indicates that the parameter is located in *http.Request.
	InRequest Location = "request"

	MediaTypeJSON = "application/json; charset=utf-8"
)

type Specification struct {
	Metadata   *Metadata
	Operations []*Operation
}

func (s *Specification) OldSpec() *openapi.Specification {
	old := new(openapi.Specification)

	old.Metadata = &openapi.Metadata{
		DocsPath:    s.Metadata.DocsPath,
		Title:       s.Metadata.Title,
		Version:     s.Metadata.Version,
		Description: s.Metadata.Description,
		BasePath:    s.Metadata.BasePath,
		DefaultTags: s.Metadata.DefaultTags,
	}

	buildParams := func(o *Operation) (params []*openapi.Param) {
		for _, b := range o.Request.Bindings {
			for _, p := range b.Params {
				params = append(params, &openapi.Param{
					Name:        b.Arg.Name,
					Type:        b.Arg.TypeString,
					RawType:     b.Arg.Type,
					In:          string(p.In),
					Alias:       p.Name,
					AliasType:   p.Type,
					Required:    p.Required,
					Description: p.Description,
					IsBlank:     b.IsBlank(),
				})
			}
		}
		return
	}

	for _, o := range s.Operations {
		old.Operations = append(old.Operations, &openapi.Operation{
			Name:    o.Name,
			Method:  o.Method,
			Pattern: o.Pattern,
			Request: openapi.Request{
				MediaType: o.Request.MediaType,
				BodyField: o.Request.BodyField,
				Params:    buildParams(o),
			},
			SuccessResponse: &openapi.Response{
				StatusCode: o.SuccessResponse.StatusCode,
				MediaType:  o.SuccessResponse.MediaType,
				Schema:     o.SuccessResponse.Schema,
				BodyField:  o.SuccessResponse.BodyField,
			},
			//FailureResponses: nil,
			Description: o.Description,
			Tags:        o.Tags,
		})
	}

	return old
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

// Parameter represents a request parameter.
// See https://spec.openapis.org/oas/v3.1.0#parameter-object.
type Parameter struct {
	In          Location // The location of the parameter.
	Name        string   // The name of the parameter.
	Required    bool     // Whether this parameter is mandatory.
	Type        string   // The type of the parameter.
	Description string   // A brief description of the parameter.
}

func (p *Parameter) UniqueKey() string {
	return fmt.Sprintf("%s.%s", p.In, p.Name)
}

// Binding represents a binding from a method argument to one or more request parameters.
type Binding struct {
	Arg    *ifacetool.Param // The method argument.
	Params []*Parameter     // The corresponding request parameters
}

// IsBlank determines whether this binding is a blank identifier.
func (b *Binding) IsBlank() bool {
	return strings.HasPrefix(b.Arg.Name, "__")
}

// IsAggregate determines whether this binding is aggregate.
func (b *Binding) IsAggregate() bool {
	return len(b.Params) > 1
}

// IsManual determines whether this binding is specified manually (i.e. by
// handwritten annotations).
func (b *Binding) IsManual() bool {
	return b.IsAggregate() || b.In() != InBody
}

func (b *Binding) In() Location {
	b.panicIfError()
	return b.Params[0].In
}

func (b *Binding) Name() string {
	b.panicIfError()
	return b.Params[0].Name
}

func (b *Binding) Required() bool {
	b.panicIfError()
	return b.Params[0].Required
}

func (b *Binding) Type() string {
	b.panicIfError()
	return b.Params[0].Type
}

func (b *Binding) Description() string {
	b.panicIfError()
	return b.Params[0].Description
}

func (b *Binding) SetIn(in Location) {
	b.panicIfError()
	b.Params[0].In = in
	if in == InPath {
		b.Params[0].Required = true
	}
}

func (b *Binding) SetName(name string) {
	b.panicIfError()
	b.Params[0].Name = name
}

func (b *Binding) SetType(typ string) {
	b.panicIfError()
	b.Params[0].Type = typ
}

func (b *Binding) SetDescription(description string) {
	b.panicIfError()
	b.Params[0].Description = description
}

func (b *Binding) panicIfError() {
	if len(b.Params) == 0 {
		panic(errors.New("invalid binding"))
	}

	if b.IsAggregate() {
		panic(errors.New("aggregate binding"))
	}
}

type Request struct {
	MediaType string

	// The name of the request field whose value is mapped to the HTTP request body.
	// Otherwise, all fields not located in path/query/header will be mapped to the HTTP body
	BodyField string

	Bindings []*Binding
}

func (r *Request) Bind(arg *ifacetool.Param, params []*Parameter) {
	b := &Binding{
		Arg:    arg,
		Params: params,
	}
	r.Bindings = append(r.Bindings, b)
}

func (r *Request) GetBinding(argName string) *Binding {
	for _, b := range r.Bindings {
		if b.Arg.Name == argName {
			return b
		}
	}
	return nil
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
	Request          *Request
	SuccessResponse  *Response
	FailureResponses []*Response
	Description      string
	Tags             []string
}

func NewOperation(name, description string) *Operation {
	op := &Operation{
		Name:        name,
		Description: description,
		Request:     new(Request),
	}
	op.Resp(http.StatusOK, MediaTypeJSON, nil)
	return op
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
