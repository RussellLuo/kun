// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package messaging

import (
	"github.com/RussellLuo/kun/pkg/oas2"
)

var (
	base = `swagger: "2.0"
info:
  title: "No Title"
  version: "0.0.0"
  description: ""
  license:
    name: "MIT"
host: "example.com"
basePath: "/"
schemes:
  - "https"
consumes:
  - "application/json"
produces:
  - "application/json"
`

	paths = `
paths:
  /messages/{messageID}:
    get:
      description: ""
      operationId: "GetMessage"
      parameters:
        - name: messageID
          in: path
          required: true
          type: string
          description: ""
      %s
  /users/{userID}/messages/{messageID}:
    get:
      description: ""
      operationId: "GetMessage1"
      parameters:
        - name: userID
          in: path
          required: true
          type: string
          description: ""
        - name: messageID
          in: path
          required: true
          type: string
          description: ""
      %s
`
)

func getResponses(schema oas2.Schema) []oas2.OASResponses {
	return []oas2.OASResponses{
		oas2.GetOASResponses(schema, "GetMessage", 200, &GetMessageResponse{}),
		oas2.GetOASResponses(schema, "GetMessage", 200, &GetMessageResponse{}),
	}
}

func getDefinitions(schema oas2.Schema) map[string]oas2.Definition {
	defs := make(map[string]oas2.Definition)

	oas2.AddResponseDefinitions(defs, schema, "GetMessage", 200, (&GetMessageResponse{}).Body())

	oas2.AddResponseDefinitions(defs, schema, "GetMessage", 200, (&GetMessageResponse{}).Body())

	return defs
}

func OASv2APIDoc(schema oas2.Schema) string {
	resps := getResponses(schema)
	paths := oas2.GenPaths(resps, paths)

	defs := getDefinitions(schema)
	definitions := oas2.GenDefinitions(defs)

	return base + paths + definitions
}
