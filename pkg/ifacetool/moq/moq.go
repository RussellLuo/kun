package moq

import (
	"go/types"

	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/ifacetool/moq/registry"
	"github.com/RussellLuo/kun/pkg/ifacetool/moq/template"
)

type Parser struct {
	cfg      Config
	registry *registry.Registry
}

type Config struct {
	SrcDir     string
	PkgName    string
	SkipEnsure bool
}

func New(cfg Config) (*Parser, error) {
	reg, err := registry.New(cfg.SrcDir, cfg.PkgName)
	if err != nil {
		return nil, err
	}
	return &Parser{
		cfg:      cfg,
		registry: reg,
	}, nil
}

func (p *Parser) Parse(ifaceName string) (*ifacetool.Data, error) {
	iface, err := p.registry.LookupInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	methods := make([]template.MethodData, iface.NumMethods())
	for j := 0; j < iface.NumMethods(); j++ {
		methods[j] = p.methodData(iface.Method(j))
	}

	data := template.Data{
		PkgName: p.mockPkgName(),
	}
	srcPkgName := p.registry.SrcPkgName()
	if srcPkgName != p.mockPkgName() {
		data.SrcPkgQualifier = p.registry.SrcPkgName() + "."
		if !p.cfg.SkipEnsure {
			imprt := p.registry.AddImport(p.registry.SrcPkg())
			data.SrcPkgQualifier = imprt.Qualifier() + "."
		}
	}

	data.Imports = p.registry.Imports()

	return toIfaceToolData(data, srcPkgName, ifaceName, methods), nil
}

func (p *Parser) methodData(f *types.Func) template.MethodData {
	sig := f.Type().(*types.Signature)

	scope := p.registry.MethodScope()
	n := sig.Params().Len()
	params := make([]template.ParamData, n)
	for i := 0; i < n; i++ {
		p := template.ParamData{
			Var: scope.AddVar(sig.Params().At(i), ""),
		}
		p.Variadic = sig.Variadic() && i == n-1 && p.Var.IsSlice() // check for final variadic argument

		params[i] = p
	}

	n = sig.Results().Len()
	results := make([]template.ParamData, n)
	for i := 0; i < n; i++ {
		results[i] = template.ParamData{
			Var: scope.AddVar(sig.Results().At(i), ""),
		}
	}

	return template.MethodData{
		Name:    f.Name(),
		Params:  params,
		Returns: results,
	}
}

func (p *Parser) mockPkgName() string {
	if p.cfg.PkgName != "" {
		return p.cfg.PkgName
	}

	return p.registry.SrcPkgName()
}

func toIfaceToolData(in template.Data, srcPkgName, ifaceName string, methods []template.MethodData) *ifacetool.Data {
	out := &ifacetool.Data{
		PkgName:         in.PkgName,
		SrcPkgName:      srcPkgName,
		SrcPkgQualifier: in.SrcPkgQualifier,
		InterfaceName:   ifaceName,
	}

	for _, imprt := range in.Imports {
		out.Imports = append(out.Imports, &ifacetool.Import{
			Alias: imprt.Alias,
			Path:  imprt.Path(),
		})
	}

	for _, meth := range methods {
		out.Methods = append(out.Methods, &ifacetool.Method{
			Name:    meth.Name,
			Params:  toIfaceToolParam(meth.Params),
			Returns: toIfaceToolParam(meth.Returns),
		})
	}

	return out
}

func toIfaceToolParam(in []template.ParamData) (out []*ifacetool.Param) {
	for _, p := range in {
		out = append(out, &ifacetool.Param{
			Name:       p.Name(),
			TypeString: p.TypeString(),
			Type:       p.Var.Type(),
			Variadic:   p.Variadic,
		})
	}
	return
}
