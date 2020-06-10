package reflector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

func GetInterfaceMethodDoc(filename, name string) (map[string][]string, error) {
	ifType, err := getAstInterfaceType(filename, name)
	if err != nil {
		return nil, err
	}

	doc := make(map[string][]string)

	for _, method := range ifType.Methods.List {
		methodName := method.Names[0].Name

		if method.Doc == nil {
			return nil, fmt.Errorf("method %q has no comment", methodName)
		}

		var comments []string
		for _, c := range method.Doc.List {
			comments = append(comments, c.Text)
		}
		doc[methodName] = comments
	}

	return doc, nil
}

func getAstInterfaceType(filename, name string) (*ast.InterfaceType, error) {
	filename, _ = filepath.Abs(filename)

	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	for _, d := range f.Decls {
		for _, s := range d.(*ast.GenDecl).Specs {
			ts, ok := s.(*ast.TypeSpec)
			if ok && ts.Name.Name == name {
				ifType, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("%q is not an interface", name)
				}
				return ifType, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find interface %q", name)
}
