package kok

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RussellLuo/kok/kok/endpoint"
	"github.com/RussellLuo/kok/kok/http/chi"
	"github.com/RussellLuo/kok/kok/http/httptest"
	"github.com/RussellLuo/kok/pkg/openapi"
	"github.com/RussellLuo/kok/pkg/reflector"
)

type Options struct {
	SchemaPtr         bool
	SchemaTag         string
	TagKeyToSnakeCase bool
	Formatted         bool
}

type Content struct {
	Endpoint []byte
	HTTP     []byte
	HTTPTest []byte
}

type Generator struct {
	endpoint *endpoint.Generator
	chi      *chi.Generator
	httptest *httptest.Generator
}

func New(opts Options) *Generator {
	return &Generator{
		endpoint: endpoint.New(endpoint.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
		}),
		chi: chi.New(chi.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
		}),
		httptest: httptest.New(opts.Formatted),
	}
}

func (g *Generator) Generate(srcFilename, interfaceName, dstPkgName, testFilename string) (content Content, err error) {
	result, err := reflector.ReflectInterface(filepath.Dir(srcFilename), dstPkgName, interfaceName)
	if err != nil {
		return content, err
	}

	// Generate the endpoint code.
	content.Endpoint, err = g.endpoint.Generate(result)
	if err != nil {
		return content, err
	}

	// Generate the HTTP code.
	doc, err := reflector.GetInterfaceMethodDoc(srcFilename, interfaceName)
	if err != nil {
		return content, err
	}

	spec, err := openapi.FromDoc(result, doc)
	if err != nil {
		return content, err
	}

	content.HTTP, err = g.chi.Generate(result, spec)
	if err != nil {
		return content, err
	}

	// Generate the HTTP tests code.
	content.HTTPTest, err = g.httptest.Generate(result, testFilename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("WARNING: Skip generating the HTTP tests due to an error (%v)\n", err)
			return content, nil
		}
		return content, err
	}

	return content, nil
}
