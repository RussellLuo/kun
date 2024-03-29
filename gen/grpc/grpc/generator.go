package grpc

import (
	"github.com/RussellLuo/kun/gen/grpc/parser"
	"github.com/RussellLuo/kun/gen/util/annotation"
	"github.com/RussellLuo/kun/gen/util/generator"
	"github.com/RussellLuo/kun/pkg/caseconv"
	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/pkgtool"
)

var (
	template = annotation.FileHeader + `
package {{.PkgInfo.CurrentPkgName}}

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"{{.PBPkgPath}}"
)

{{- $pbPkgPrefix := .PBPkgPrefix}}
{{- $endpointPkgPrefix := .PkgInfo.EndpointPkgPrefix}}
{{- $serviceName := .Data.InterfaceName}}

type grpcServer struct {
	{{$pbPkgPrefix}}Unimplemented{{$serviceName}}Server

	{{range .Service.RPCs -}}
	{{lowerFirst .Name}} kitgrpc.Handler
	{{end -}} {{/* range .Service.RPCs */}}
}

{{- range .Service.RPCs}}

func (s *grpcServer) {{.Name}}(ctx context.Context, req *{{$pbPkgPrefix}}{{.Request.Name}}) (*{{$pbPkgPrefix}}{{.Response.Name}}, error) {
	_, resp, err := s.{{lowerFirst .Name}}.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*{{$pbPkgPrefix}}{{.Response.Name}}), nil
}
{{- end}} {{/* range .Service.RPCs */}}

func NewGRPCServer(svc {{$.Data.SrcPkgQualifier}}{{$serviceName}}, codecs grpccodec.Codecs) {{$pbPkgPrefix}}{{$serviceName}}Server {
	var codec grpccodec.Codec
	s := new(grpcServer)

	{{- range .Service.RPCs}}

	codec = codecs.EncodeDecoder("{{.Name}}")
	s.{{lowerFirst .Name}} = kitgrpc.NewServer(
		{{$endpointPkgPrefix}}MakeEndpointOf{{.Name}}(svc),
		decode{{.Request.Name}}(codec),
		encode{{.Response.Name}}(codec),
	)
	{{- end}} {{/* range .Service.RPCs */}}

	return s
}

{{- range .Service.RPCs}}

// decode{{.Request.Name}} converts a gRPC request to an endpoint request.
func decode{{.Request.Name}}(codec grpccodec.Codec) kitgrpc.DecodeRequestFunc {
	return func(_ context.Context, grpcReq interface{}) (interface{}, error) {
		var req {{$endpointPkgPrefix}}{{.Request.Name}}
		pb := grpcReq.(*{{$pbPkgPrefix}}{{.Request.Name}})
		if err := codec.DecodeRequest(pb, &req); err != nil {
			return nil, err
		}
		return {{ampersand}}req, nil
	}
}

// encode{{.Response.Name}} converts an endpoint response to a gRPC response.
func encode{{.Response.Name}}(codec grpccodec.Codec) kitgrpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		pb := new({{$pbPkgPrefix}}{{.Response.Name}})
		resp := response.({{asterisks}}{{$endpointPkgPrefix}}{{.Response.Name}})
		if err := codec.EncodeResponse(resp, pb); err != nil {
			return nil, err
		}
		return pb, nil
	}
}
{{- end}} {{/* range .Service.RPCs */}}
`
)

type Options struct {
	SchemaPtr bool
	SchemaTag string
	Formatted bool
}

type Generator struct {
	opts *Options
}

func New(opts *Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(pkgInfo *generator.PkgInfo, pbOutDir string, ifaceData *ifacetool.Data, service *parser.Service) (*generator.File, error) {
	data := struct {
		PBPkgPath   string
		PBPkgPrefix string
		Data        *ifacetool.Data
		PkgInfo     *generator.PkgInfo
		Service     *parser.Service
	}{
		PBPkgPath:   pkgtool.PkgPathFromDir(pbOutDir),
		PBPkgPrefix: pkgtool.PkgNameFromDir(pbOutDir) + ".",
		Data:        ifaceData,
		PkgInfo:     pkgInfo,
		Service:     service,
	}

	return generator.Generate(template, data, generator.Options{
		Funcs: map[string]interface{}{
			"ampersand": func() string {
				if g.opts.SchemaPtr {
					return "&"
				}
				return ""
			},
			"asterisks": func() string {
				if g.opts.SchemaPtr {
					return "*"
				}
				return ""
			},
			"lowerFirst": caseconv.LowerFirst,
		},
		Formatted:      g.opts.Formatted,
		TargetFileName: "grpc.go",
	})
}
