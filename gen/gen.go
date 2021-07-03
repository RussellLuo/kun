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
	OutDir        string
	FlatLayout    bool
	SchemaPtr     bool
	SchemaTag     string
	Formatted     bool
	SnakeCase     bool
	EnableTracing bool
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

func (g *Generator) Generate(srcFilename, interfaceName, testFilename string) (files []*generator.File, err error) {
	result, err := g.getInterfaceResult(srcFilename, interfaceName)
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
	epFile, err := g.generateEndpoint(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, epFile)

	// Generate the HTTP code.
	httpFiles, err := g.generateHTTP(result, spec, testFilename)
	if err != nil {
		return files, err
	}
	files = append(files, httpFiles...)

	return files, nil
}

func (g *Generator) generateEndpoint(result *reflector.Result, spec *openapi.Specification) (file *generator.File, err error) {
	outDir := g.getOutDir("endpoint")
	if err = ensureDir(outDir); err != nil {
		return
	}
	defer func() {
		if file != nil {
			moveTo(outDir, file)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)
	file, err = g.endpoint.Generate(pkgInfo, result, spec)
	if err != nil {
		return
	}

	return
}

func (g *Generator) generateHTTP(result *reflector.Result, spec *openapi.Specification, testFilename string) (files []*generator.File, err error) {
	outDir := g.getOutDir("http")
	if err := ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			moveTo(outDir, f)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)

	// Generate the HTTP server code.
	f, err := g.chi.Generate(pkgInfo, result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	// Generate the HTTP client code.
	f, err = g.httpclient.Generate(pkgInfo, result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	// Generate the HTTP tests code.
	f, err = g.httptest.Generate(pkgInfo, result, testFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			return files, err
		}
		fmt.Printf("WARNING: Skip generating the HTTP tests due to an error (%v)\n", err)
	}
	if f != nil {
		files = append(files, f)
	}

	// Generate the helper OASv2 code.
	f, err = g.oasv2.Generate(pkgInfo, result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

func (g *Generator) getInterfaceResult(srcFilename, interfaceName string) (*reflector.Result, error) {
	pkgName := ""
	if !g.opts.FlatLayout {
		// Non-empty pkgName makes all type names used in the interface full-qualified.
		pkgName = "x"
	}

	result, err := reflector.ReflectInterface(filepath.Dir(srcFilename), pkgName, interfaceName)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (g *Generator) getOutDir(sub string) string {
	dir := g.opts.OutDir
	if !g.opts.FlatLayout {
		dir = filepath.Join(dir, sub)
	}
	return dir
}

func (g *Generator) getPkgInfo(dir string) *generator.PkgInfo {
	pkgInfo := &generator.PkgInfo{
		CurrentPkgName: reflector.PkgNameFromDir(dir),
	}
	if !g.opts.FlatLayout {
		pkgInfo.EndpointPkgPrefix = reflector.PkgNameFromDir(g.getOutDir("endpoint")) + "."
		pkgInfo.EndpointPkgPath = reflector.PkgPathFromDir(g.getOutDir("endpoint"))
	}
	return pkgInfo
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func moveTo(dir string, f *generator.File) {
	f.Name = filepath.Join(dir, f.Name)
}
