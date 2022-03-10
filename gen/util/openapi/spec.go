package openapi

import (
	"go/types"
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
	GoMethodName     string
	Method           string
	Pattern          string
	Request          Request
	SuccessResponse  *Response
	FailureResponses []*Response
	Description      string
	Tags             []string
}
