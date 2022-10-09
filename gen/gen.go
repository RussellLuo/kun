package gen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	crongenerator "github.com/RussellLuo/kun/gen/cron/generator"
	cronparser "github.com/RussellLuo/kun/gen/cron/parser"
	"github.com/RussellLuo/kun/gen/endpoint"
	eventgenerator "github.com/RussellLuo/kun/gen/event/generator"
	eventparser "github.com/RussellLuo/kun/gen/event/parser"
	"github.com/RussellLuo/kun/gen/grpc/grpc"
	grpcparser "github.com/RussellLuo/kun/gen/grpc/parser"
	"github.com/RussellLuo/kun/gen/grpc/proto"
	"github.com/RussellLuo/kun/gen/http/chi"
	"github.com/RussellLuo/kun/gen/http/httpclient"
	"github.com/RussellLuo/kun/gen/http/oas2"
	httpparser "github.com/RussellLuo/kun/gen/http/parser"
	"github.com/RussellLuo/kun/gen/util/docutil"
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
}

type Generator struct {
	endpoint   *endpoint.Generator
	chi        *chi.Generator
	httpclient *httpclient.Generator
	oas2       *oas2.Generator
	proto      *proto.Generator
	grpc       *grpc.Generator
	event      *eventgenerator.Generator
	cron       *crongenerator.Generator

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
		event: eventgenerator.New(&eventgenerator.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
		}),
		cron: crongenerator.New(&crongenerator.Options{
			SchemaPtr: opts.SchemaPtr,
			SchemaTag: opts.SchemaTag,
			Formatted: opts.Formatted,
		}),
		opts: opts,
	}
}

func (g *Generator) Generate(srcFilename, interfaceName string) (files []*generator.File, err error) {
	data, err := g.parseInterface(srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	newSpec, transports, err := httpparser.Parse(data, g.opts.SnakeCase)
	if err != nil {
		return nil, err
	}
	spec := newSpec.OldSpec()

	epFile, err := g.generateEndpoint(data, spec)
	if err != nil {
		return files, err
	}
	files = append(files, epFile)

	switch mergeTransports(transports) {
	case docutil.TransportHTTP:
		httpFiles, err := g.generateHTTP(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

	case docutil.TransportGRPC:
		grpcFiles, err := g.generateGRPC(data)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)

	case docutil.TransportEvent:
		eventFiles, err := g.generateEvent(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, eventFiles...)

	case docutil.TransportCron:
		cronFiles, err := g.generateCron(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, cronFiles...)

	case docutil.TransportAll:
		httpFiles, err := g.generateHTTP(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, httpFiles...)

		grpcFiles, err := g.generateGRPC(data)
		if err != nil {
			return files, err
		}
		files = append(files, grpcFiles...)

		eventFiles, err := g.generateEvent(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, eventFiles...)

		cronFiles, err := g.generateCron(data, spec)
		if err != nil {
			return files, err
		}
		files = append(files, cronFiles...)
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
func (g *Generator) generateHTTP(data *ifacetool.Data, spec *openapi.Specification) (files []*generator.File, err error) {
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

	// Generate the helper OAS2 code.
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
	f, err := g.proto.Generate(pbOutDir, data, service)
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
	f, err = g.grpc.Generate(pkgInfo, pbOutDir, data, service)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

// generateEvent generates the event code.
func (g *Generator) generateEvent(data *ifacetool.Data, spec *openapi.Specification) (files []*generator.File, err error) {
	outDir := g.getOutDir("event")
	if err := ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			f.MoveTo(outDir)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)

	eventInfo, err := eventparser.Parse(data, g.opts.SnakeCase)
	if err != nil {
		return files, err
	}

	f, err := g.event.Generate(pkgInfo, data, eventInfo, spec)
	if err != nil {
		return files, err
	}
	files = append(files, f)

	return files, nil
}

// generateCron generates the cron code.
func (g *Generator) generateCron(data *ifacetool.Data, spec *openapi.Specification) (files []*generator.File, err error) {
	outDir := g.getOutDir("cron")
	if err := ensureDir(outDir); err != nil {
		return files, err
	}
	defer func() {
		for _, f := range files {
			f.MoveTo(outDir)
		}
	}()

	pkgInfo := g.getPkgInfo(outDir)

	cronSpec, err := cronparser.Parse(data, g.opts.SnakeCase)
	if err != nil {
		return files, err
	}

	f, err := g.cron.Generate(pkgInfo, data, cronSpec, spec)
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

func mergeTransports(transports []docutil.Transport) (result docutil.Transport) {
	for _, t := range transports {
		result = result | t
	}
	return result
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
