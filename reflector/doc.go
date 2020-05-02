package reflector

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
)

type MethodDoc struct {
	Name     string
	Comments []string
}

func GetInterfaceDoc(filename, objName string) (docs []MethodDoc) {
	ifType, err := getInterfaceType(filename, objName)
	if err != nil {
		panic(err)
	}

	for _, field := range ifType.Methods.List {
		d := MethodDoc{
			Name: field.Names[0].Name,
		}
		for _, c := range field.Doc.List {
			d.Comments = append(d.Comments, c.Text)
		}

		docs = append(docs, d)
	}

	return
}

func getInterfaceType(filename, objName string) (*ast.InterfaceType, error) {
	filename, _ = filepath.Abs(filename)

	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	docPkg := doc.New(&ast.Package{
		Files: map[string]*ast.File{filename: f},
	}, "", doc.AllDecls)

	for _, t := range docPkg.Types {
		for _, s := range t.Decl.Specs {
			ts := s.(*ast.TypeSpec)
			if ts.Name.Name == objName {
				ifType, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("%s(ts.Type) is not an interface", objName)
				}
				return ifType, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find interface %s", objName)
}
