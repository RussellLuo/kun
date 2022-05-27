package pkgtool

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/RussellLuo/kun/pkg/ifacetool"
	"github.com/RussellLuo/kun/pkg/ifacetool/moq"
	"golang.org/x/tools/go/packages"
)

func ParseInterface(pkgName, srcFilename, interfaceName string) (*ifacetool.Data, error) {
	moqParser, err := moq.New(moq.Config{
		SrcDir:  filepath.Dir(srcFilename),
		PkgName: pkgName,
	})
	if err != nil {
		return nil, err
	}

	data, err := moqParser.Parse(interfaceName)
	if err != nil {
		return nil, err
	}

	doc, err := newInterfaceDoc(srcFilename, interfaceName)
	if err != nil {
		return nil, err
	}

	data.InterfaceDoc = doc.Doc
	for _, m := range data.Methods {
		m.Doc = doc.MethodDocs[m.Name]
	}

	return data, nil
}

type interfaceDoc struct {
	Doc        []string
	MethodDocs map[string][]string
}

func newInterfaceDoc(filename, name string) (*interfaceDoc, error) {
	ifType, ifDoc, err := getAstInterfaceInfo(filename, name)
	if err != nil {
		return nil, err
	}

	var doc []string
	if ifDoc != nil {
		for _, c := range ifDoc.List {
			doc = append(doc, c.Text)
		}
	}

	methodDocs := make(map[string][]string)

	for _, method := range ifType.Methods.List {
		methodName := method.Names[0].Name

		if method.Doc == nil {
			continue
		}

		var comments []string
		for _, c := range method.Doc.List {
			comments = append(comments, c.Text)
		}
		methodDocs[methodName] = comments
	}

	return &interfaceDoc{Doc: doc, MethodDocs: methodDocs}, nil
}

func getAstInterfaceInfo(filename, name string) (*ast.InterfaceType, *ast.CommentGroup, error) {
	filename, _ = filepath.Abs(filename)

	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return nil, nil, err
	}

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			// Interface definition is a type declaration, so we ignore
			// non-GenDecl nodes here.
			continue
		}

		for _, s := range gd.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if ok && ts.Name.Name == name {
				doc := getInterfaceDoc(ts, gd)
				ifType, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, nil, fmt.Errorf("%q is not an interface", name)
				}
				return ifType, doc, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("could not find interface %q", name)
}

func getInterfaceDoc(ts *ast.TypeSpec, gd *ast.GenDecl) *ast.CommentGroup {
	// See https://github.com/golang/go/issues/27477

	if ts.Doc != nil {
		// Use the documentation of ts (an individual type declaration), if any.
		return ts.Doc
	}

	if len(gd.Specs) == 1 && gd.Doc != nil {
		// If gd (a grouped type declaration) only has one inner type declaration,
		// use the documentation of gd, if any.
		return gd.Doc
	}

	return nil
}

func PkgPathFromDir(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	pkg, err := pkgInfoFromPath(abs, packages.NeedModule)
	if err != nil {
		panic(err)
	}

	if pkg == nil || pkg.Module == nil {
		return ""
	}

	// Remove the module root directory.
	rel, err := filepath.Rel(pkg.Module.Dir, abs)
	if err != nil {
		panic(err)
	}
	// Add the module path prefix.
	modPath := filepath.Join(pkg.Module.Path, rel)

	// The final module path must be separated by slashes ('/').
	return filepath.ToSlash(modPath)
}

func PkgNameFromDir(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	pkg, err := pkgInfoFromPath(abs, packages.NeedName)
	if err != nil {
		panic(err)
	}
	if pkg != nil && pkg.Name != "" {
		return pkg.Name
	}

	// Default to the directory name.
	return filepath.Base(abs)
}

func pkgInfoFromPath(srcDir string, mode packages.LoadMode) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: mode,
		Dir:  srcDir,
	})
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("no packages found")
	}
	if len(pkgs) > 1 {
		return nil, errors.New("more than one package was found")
	}
	return pkgs[0], nil
}
