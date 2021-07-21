package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RussellLuo/kok/gen/endpoint"
	"github.com/RussellLuo/kok/gen/grpc/grpc"
	"github.com/RussellLuo/kok/gen/grpc/parser"
	"github.com/RussellLuo/kok/gen/grpc/proto"
	"github.com/RussellLuo/kok/gen/http/chi"
	"github.com/RussellLuo/kok/gen/http/httpclient"
	"github.com/RussellLuo/kok/gen/http/httptest"
	"github.com/RussellLuo/kok/gen/http/oasv2"
	"github.com/RussellLuo/kok/gen/util/generator"
	"github.com/RussellLuo/kok/gen/util/openapi"
	"github.com/RussellLuo/kok/gen/util/reflector"
	"github.com/RussellLuo/kok/pkg/ifacetool"
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
	proto      *proto.Generator
	grpc       *grpc.Generator

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
		proto: proto.New(&proto.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
		}),
		grpc: grpc.New(&grpc.Options{
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

	spec, transports, err := openapi.FromDoc(result, doc, g.opts.SnakeCase)
	if err != nil {
		return nil, err
	}

	epFile, err := g.generateEndpoint(result, spec)
	if err != nil {
		return files, err
	}
	files = append(files, epFile)

	switch mergeTransports(transports) {
	case openapi.TransportHTTP:
		httpFiles, err := g.generateHTTP(result.Data, spec, testFilename)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

	case openapi.TransportGRPC:
		grpcFiles, err := g.generateGRPC(result, doc)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)

	case openapi.TransportAll:
		httpFiles, err := g.generateHTTP(result.Data, spec, testFilename)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

		grpcFiles, err := g.generateGRPC(result, doc)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)
	}

	return files, nil
}

// generateEndpoint generates the endpoint code.
func (g *Generator) generateEndpoint(result *reflector.Result, spec *openapi.Specification) (file *generator.File, err error) {
	outDir := g.getOutDir("endpoint")
	if err = ensureDir(outDir); err != nil {
		return
	}
	defer func() {
		if file != nil {
			file.MoveTo(outDir)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)
	file, err = g.endpoint.Generate(pkgInfo, result, spec)
	if err != nil {
		return
	}

	return
}

// generateHTTP generates the HTTP code.
func (g *Generator) generateHTTP(data *ifacetool.Data, spec *openapi.Specification, testFilename string) (files []*generator.File, err error) {
	outDir := g.getOutDir("http")
	if err := ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			f.MoveTo(outDir)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)

	// Generate the HTTP server code.
	f, err := g.chi.Generate(pkgInfo, data, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	// Generate the HTTP client code.
	f, err = g.httpclient.Generate(pkgInfo, data, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	// Generate the HTTP tests code.
	f, err = g.httptest.Generate(pkgInfo, data, testFilename)
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
	f, err = g.oasv2.Generate(pkgInfo, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

// generateGRPC generates the gRPC code.
func (g *Generator) generateGRPC(result *reflector.Result, doc *reflector.InterfaceDoc) (files []*generator.File, err error) {
	outDir := g.getOutDir("grpc")
	if err = ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			f.MoveTo(outDir)
		}
	}()

	service, err := parser.Parse(result.Data, doc)
	if err != nil {
		return files, err
	}

	// Generate the `.proto` file.
	pbOutDir := filepath.Join(outDir, "pb")
	if err = ensureDir(pbOutDir); err != nil {
		return files, err
	}
	pbPkgPath := reflector.PkgPathFromDir(pbOutDir)
	f, err := g.proto.Generate(pbPkgPath, result, service)
	if err != nil {
		return files, err
	}
	// Write the `proto` file at once.
	f.MoveTo(pbOutDir)
	if err := f.Write(); err != nil {
		return files, err
	}

	// Compile the `.proto` file to the gRPC definition.
	// See https://grpc.io/docs/languages/go/basics/#generating-client-and-server-code
	cmd := exec.Command("protoc",
		"--go_out=.", "--go_opt=paths=source_relative",
		"--go-grpc_out=.", "--go-grpc_opt=paths=source_relative",
		filepath.Join(pbOutDir, result.SrcPkgName+".proto"),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return files, fmt.Errorf("failed to compile proto: %s", out)
	}

	// Generate the glue code for adapting the gRPC definition to Go kit.
	pkgInfo := g.getPkgInfo(outDir)
	f, err = g.grpc.Generate(pkgInfo, pbPkgPath, result, service)
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

func mergeTransports(transports []openapi.Transport) (result openapi.Transport) {
	for _, t := range transports {
		result = result | t
	}
	return result
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
