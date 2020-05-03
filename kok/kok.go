package kok

import (
	"path/filepath"

	"github.com/RussellLuo/kok/kok/endpoint"
	"github.com/RussellLuo/kok/kok/http"
	"github.com/RussellLuo/kok/oapi"
	"github.com/RussellLuo/kok/reflector"
)

type Options struct {
	SchemaPtr         bool
	SchemaTag         string
	TagKeyToSnakeCase bool
}

type Content struct {
	Endpoint []byte
	HTTP     []byte
}

type Generator struct {
	opts     Options
	endpoint *endpoint.Generator
	chi      *http.ChiGenerator
}

func New(opts Options) *Generator {
	return &Generator{
		opts: opts,
		endpoint: endpoint.New(endpoint.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
		}),
		chi: http.NewChi(http.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
		}),
	}
}

func (g *Generator) Generate(srcFilename, interfaceName, dstPkgName string) (content Content, err error) {
	result, err := reflector.ReflectInterface(filepath.Dir(srcFilename), dstPkgName, interfaceName)
	if err != nil {
		return content, err
	}

	content.Endpoint, err = g.endpoint.Generate(result)
	if err != nil {
		return content, err
	}

	doc, err := reflector.GetInterfaceMethodDoc(srcFilename, interfaceName)
	if err != nil {
		return content, err
	}

	spec, err := oapi.FromDoc(result, doc)
	if err != nil {
		return content, err
	}

	content.HTTP, err = g.chi.Generate(result, spec)
	if err != nil {
		return content, err
	}

	return content, nil
}
