package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RussellLuo/kun/gen/endpoint"
	"github.com/RussellLuo/kun/gen/grpc/grpc"
	grpcparser "github.com/RussellLuo/kun/gen/grpc/parser"
	"github.com/RussellLuo/kun/gen/grpc/proto"
	"github.com/RussellLuo/kun/gen/http/chi"
	"github.com/RussellLuo/kun/gen/http/httpclient"
	"github.com/RussellLuo/kun/gen/http/httptest"
	"github.com/RussellLuo/kun/gen/http/oas2"
	httpparser "github.com/RussellLuo/kun/gen/http/parser"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/gen/util/openapi"
	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/pkgtool"
)

type Options struct {
	OutDir        string
	FlatLayout    bool
	SchemaPtr     bool
	SchemaTag     string
	Formatted     bool
	SnakeCase     bool
	EnableTracing bool
	OldAnnotation bool
}

type Generator struct {
	endpoint   *endpoint.Generator
	chi        *chi.Generator
	httptest   *httptest.Generator
	httpclient *httpclient.Generator
	oas2       *oas2.Generator
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
		oas2: oas2.New(&oas2.Options{
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
	data, err := g.parseInterface(srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	var (
		spec       *openapi.Specification
		transports []openapi.Transport
	)
	if g.opts.OldAnnotation {
		spec, transports, err = openapi.FromDoc(data, g.opts.SnakeCase)
		if err != nil {
			return nil, err
		}
	} else {
		newSpec, newTransports, err := httpparser.Parse(data, g.opts.SnakeCase)
		if err != nil {
			return nil, err
		}
		spec = newSpec.OldSpec()
		for _, t := range newTransports {
			transports = append(transports, openapi.Transport(t))
		}
	}

	epFile, err := g.generateEndpoint(data, spec)
	if err != nil {
		return files, err
	}
	files = append(files, epFile)

	switch mergeTransports(transports) {
	case openapi.TransportHTTP:
		httpFiles, err := g.generateHTTP(data, spec, testFilename)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

	case openapi.TransportGRPC:
		grpcFiles, err := g.generateGRPC(data)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)

	case openapi.TransportAll:
		httpFiles, err := g.generateHTTP(data, spec, testFilename)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

		grpcFiles, err := g.generateGRPC(data)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)
	}

	return files, nil
}

// generateEndpoint generates the endpoint code.
func (g *Generator) generateEndpoint(data *ifacetool.Data, spec *openapi.Specification) (file *generator.File, err error) {
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
	file, err = g.endpoint.Generate(pkgInfo, data, spec)
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
	f, err = g.oas2.Generate(pkgInfo, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

// generateGRPC generates the gRPC code.
func (g *Generator) generateGRPC(data *ifacetool.Data) (files []*generator.File, err error) {
	outDir := g.getOutDir("grpc")
	if err = ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			f.MoveTo(outDir)
		}
	}()

	service, err := grpcparser.Parse(data)
	if err != nil {
		return files, err
	}

	// Generate the `.proto` file.
	pbOutDir := filepath.Join(outDir, "pb")
	if err = ensureDir(pbOutDir); err != nil {
		return files, err
	}
	pbPkgPath := pkgtool.PkgPathFromDir(pbOutDir)
	f, err := g.proto.Generate(pbPkgPath, data, service)
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
		filepath.Join(pbOutDir, data.SrcPkgName+".proto"),
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return files, fmt.Errorf("failed to compile proto: %s", out)
	}

	// Generate the glue code for adapting the gRPC definition to Go kit.
	pkgInfo := g.getPkgInfo(outDir)
	f, err = g.grpc.Generate(pkgInfo, pbPkgPath, data, service)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

func (g *Generator) parseInterface(srcFilename, interfaceName string) (*ifacetool.Data, error) {
	pkgName := ""
	if !g.opts.FlatLayout {
		// Non-empty pkgName makes all type names used in the interface full-qualified.
		pkgName = "x"
	}

	data, err := pkgtool.ParseInterface(pkgName, srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	return data, nil
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
		CurrentPkgName: pkgtool.PkgNameFromDir(dir),
	}
	if !g.opts.FlatLayout {
		pkgInfo.EndpointPkgPrefix = pkgtool.PkgNameFromDir(g.getOutDir("endpoint")) + "."
		pkgInfo.EndpointPkgPath = pkgtool.PkgPathFromDir(g.getOutDir("endpoint"))
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
