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
	"github.com/RussellLuo/kok/gen/util/generator"
	"github.com/RussellLuo/kok/gen/util/openapi"
	"github.com/RussellLuo/kok/gen/util/reflector"
)

type Options struct {
	SchemaPtr     bool
	SchemaTag     string
	Formatted     bool
	SnakeCase     bool
	EnableTracing bool
	OutDir        string
}

type Generator struct {
	endpoint   *endpoint.Generator
	chi        *chi.Generator
	httptest   *httptest.Generator
	httpclient *httpclient.Generator
	oasv2      *oasv2.Generator

	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{
		endpoint: endpoint.New(&endpoint.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
			SnakeCase: opts.SnakeCase,
		}),
		chi: chi.New(&chi.Options{
			SchemaPtr:     opts.SchemaPtr,
			SchemaTag:     opts.SchemaTag,
			Formatted:     opts.Formatted,
			EnableTracing: opts.EnableTracing,
		}),
		httptest: httptest.New(&httptest.Options{
			Formatted: opts.Formatted,
		}),
		httpclient: httpclient.New(&httpclient.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
		}),
		oasv2: oasv2.New(&oasv2.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
		}),
		opts: opts,
	}
}

func (g *Generator) Generate(srcFilename, interfaceName, dstPkgName, testFilename string) (files []*generator.File, err error) {
	result, err := reflector.ReflectInterface(filepath.Dir(srcFilename), dstPkgName, interfaceName)
	if err != nil {
		return nil, err
	}

	doc, err := reflector.NewInterfaceDoc(srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	spec, err := openapi.FromDoc(result, doc, g.opts.SnakeCase)
	if err != nil {
		return nil, err
	}

	// Generate the endpoint code.
	f, err := g.endpoint.Generate(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	// Generate the HTTP code.
	httpFiles, err := g.generateHTTP(result, spec, testFilename)
	if err != nil {
		return files, err
	}
	files = append(files, httpFiles...)

	return files, nil
}

func (g *Generator) generateHTTP(result *reflector.Result, spec *openapi.Specification, testFilename string) (files []*generator.File, err error) {
	outDir := g.opts.OutDir
	if err := ensureDir(outDir); err != nil {
		return files, err
	}

	// Generate the HTTP server code.
	f, err := g.chi.Generate(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, moveTo(f, outDir))

	// Generate the HTTP client code.
	f, err = g.httpclient.Generate(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, moveTo(f, outDir))

	// Generate the HTTP tests code.
	f, err = g.httptest.Generate(result, testFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			return files, err
		}
		fmt.Printf("WARNING: Skip generating the HTTP tests due to an error (%v)\n", err)
	}
	if f != nil {
		files = append(files, moveTo(f, outDir))
	}

	// Generate the helper OASv2 code.
	f, err = g.oasv2.Generate(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, moveTo(f, outDir))

	return files, nil
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func moveTo(f *generator.File, dir string) *generator.File {
	f.Name = filepath.Join(dir, f.Name)
	return f
}
