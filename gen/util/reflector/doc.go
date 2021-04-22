package reflector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type InterfaceDoc struct {
	Doc        []string
	MethodDocs map[string][]string
}

func NewInterfaceDoc(filename, name string) (*InterfaceDoc, error) {
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

	return &InterfaceDoc{Doc: doc, MethodDocs: methodDocs}, nil
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
