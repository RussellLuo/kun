// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package profilesvc

import (
	"reflect"

	"github.com/RussellLuo/kun/pkg/oas2"
)

var (
	base = `swagger: "2.0"
info:
  title: "No Title"
  version: "0.0.0"
  description: "Service is a simple CRUD interface for user profiles."
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
  /profiles/{id}/addresses/{addressID}:
    delete:
      description: ""
      operationId: "DeleteAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: addressID
          in: path
          required: true
          type: string
          description: ""
      %s
    get:
      description: ""
      operationId: "GetAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: addressID
          in: path
          required: true
          type: string
          description: ""
      %s
  /profiles/{id}:
    delete:
      description: ""
      operationId: "DeleteProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
      %s
    get:
      description: ""
      operationId: "GetProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
      %s
    patch:
      description: ""
      operationId: "PatchProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PatchProfileRequestBody"
      %s
    put:
      description: ""
      operationId: "PutProfile"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PutProfileRequestBody"
      %s
  /profiles/{id}/addresses:
    get:
      description: ""
      operationId: "GetAddresses"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
      %s
    post:
      description: ""
      operationId: "PostAddress"
      parameters:
        - name: id
          in: path
          required: true
          type: string
          description: ""
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PostAddressRequestBody"
      %s
  /profiles:
    post:
      description: ""
      operationId: "PostProfile"
      parameters:
        - name: body
          in: body
          schema:
            $ref: "#/definitions/PostProfileRequestBody"
      %s
`
)

func getResponses(schema oas2.Schema) []oas2.OASResponses {
	return []oas2.OASResponses{
		oas2.GetOASResponses(schema, "DeleteAddress", 200, &DeleteAddressResponse{}),
		oas2.GetOASResponses(schema, "GetAddress", 200, &GetAddressResponse{}),
		oas2.GetOASResponses(schema, "DeleteProfile", 200, &DeleteProfileResponse{}),
		oas2.GetOASResponses(schema, "GetProfile", 200, &GetProfileResponse{}),
		oas2.GetOASResponses(schema, "PatchProfile", 200, &PatchProfileResponse{}),
		oas2.GetOASResponses(schema, "PutProfile", 200, &PutProfileResponse{}),
		oas2.GetOASResponses(schema, "GetAddresses", 200, &GetAddressesResponse{}),
		oas2.GetOASResponses(schema, "PostAddress", 200, &PostAddressResponse{}),
		oas2.GetOASResponses(schema, "PostProfile", 200, &PostProfileResponse{}),
	}
}

func getDefinitions(schema oas2.Schema) map[string]oas2.Definition {
	defs := make(map[string]oas2.Definition)

	oas2.AddResponseDefinitions(defs, schema, "DeleteAddress", 200, (&DeleteAddressResponse{}).Body())

	oas2.AddResponseDefinitions(defs, schema, "DeleteProfile", 200, (&DeleteProfileResponse{}).Body())

	oas2.AddResponseDefinitions(defs, schema, "GetAddress", 200, (&GetAddressResponse{}).Body())

	oas2.AddResponseDefinitions(defs, schema, "GetAddresses", 200, (&GetAddressesResponse{}).Body())

	oas2.AddResponseDefinitions(defs, schema, "GetProfile", 200, (&GetProfileResponse{}).Body())

	oas2.AddDefinition(defs, "PatchProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oas2.AddResponseDefinitions(defs, schema, "PatchProfile", 200, (&PatchProfileResponse{}).Body())

	oas2.AddDefinition(defs, "PostAddressRequestBody", reflect.ValueOf(&struct {
		Address Address `json:"address"`
	}{}))
	oas2.AddResponseDefinitions(defs, schema, "PostAddress", 200, (&PostAddressResponse{}).Body())

	oas2.AddDefinition(defs, "PostProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oas2.AddResponseDefinitions(defs, schema, "PostProfile", 200, (&PostProfileResponse{}).Body())

	oas2.AddDefinition(defs, "PutProfileRequestBody", reflect.ValueOf(&struct {
		Profile Profile `json:"profile"`
	}{}))
	oas2.AddResponseDefinitions(defs, schema, "PutProfile", 200, (&PutProfileResponse{}).Body())

	return defs
}

func OASv2APIDoc(schema oas2.Schema) string {
	resps := getResponses(schema)
	paths := oas2.GenPaths(resps, paths)

	defs := getDefinitions(schema)
	definitions := oas2.GenDefinitions(defs)

	return base + paths + definitions
}
