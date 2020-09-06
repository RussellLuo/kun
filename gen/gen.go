package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RussellLuo/kok/gen/endpoint"
	"github.com/RussellLuo/kok/gen/http/chi"
	"github.com/RussellLuo/kok/gen/http/httpclient"
	"github.com/RussellLuo/kok/gen/http/httptest"
	"github.com/RussellLuo/kok/gen/http/oasv2"
	"github.com/RussellLuo/kok/gen/util/openapi"
	"github.com/RussellLuo/kok/gen/util/reflector"
)

type Options struct {
	SchemaPtr         bool
	SchemaTag         string
	TagKeyToSnakeCase bool
	Formatted         bool
	EnableTracing     bool
}

type Content struct {
	Endpoint   []byte
	HTTP       []byte
	HTTPTest   []byte
	HTTPClient []byte
	OASv2      []byte
}

type Generator struct {
	endpoint   *endpoint.Generator
	chi        *chi.Generator
	httptest   *httptest.Generator
	httpclient *httpclient.Generator
	oasv2      *oasv2.Generator
}

func New(opts Options) *Generator {
	return &Generator{
		endpoint: endpoint.New(&endpoint.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
		}),
		chi: chi.New(&chi.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
			EnableTracing:     opts.EnableTracing,
		}),
		httptest: httptest.New(&httptest.Options{
			Formatted: opts.Formatted,
		}),
		httpclient: httpclient.New(&httpclient.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
		}),
		oasv2: oasv2.New(&oasv2.Options{
			SchemaPtr:         opts.SchemaPtr,
			SchemaTag:         opts.SchemaTag,
			TagKeyToSnakeCase: opts.TagKeyToSnakeCase,
			Formatted:         opts.Formatted,
		}),
	}
}

func (g *Generator) Generate(srcFilename, interfaceName, dstPkgName, testFilename string) (content Content, err error) {
	result, err := reflector.ReflectInterface(filepath.Dir(srcFilename), dstPkgName, interfaceName)
	if err != nil {
		return content, err
	}

	spec, err := g.getSpec(result, srcFilename, interfaceName)
	if err != nil {
		return content, err
	}

	// Generate the endpoint code.
	content.Endpoint, err = g.endpoint.Generate(result, spec)
	if err != nil {
		return content, err
	}

	// Generate the HTTP code.
	content.HTTP, err = g.chi.Generate(result, spec)
	if err != nil {
		return content, err
	}

	// Generate the HTTP client code.
	content.HTTPClient, err = g.httpclient.Generate(result, spec)
	if err != nil {
		return content, err
	}

	// Generate the HTTP tests code.
	content.HTTPTest, err = g.httptest.Generate(result, testFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			return content, err
		}
		fmt.Printf("WARNING: Skip generating the HTTP tests due to an error (%v)\n", err)
	}

	content.OASv2, err = g.oasv2.Generate(result, spec)
	if err != nil {
		return content, err
	}

	return content, nil
}

func (g *Generator) getSpec(result *reflector.Result, srcFilename, interfaceName string) (*openapi.Specification, error) {
	doc, err := reflector.GetInterfaceMethodDoc(srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	spec, err := openapi.FromDoc(result, doc)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
