package reflector

import (
	"errors"
	"fmt"
	"go/build"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Most of the following code is borrowed from
// https://github.com/matryer/moq/blob/master/pkg/moq/moq.go

type Method struct {
	Name    string
	Params  []*Param
	Returns []*Param
}

func (m *Method) Arglist() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.String()
	}
	return strings.Join(params, ", ")
}

func (m *Method) ArgCallList() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.CallName()
	}
	return strings.Join(params, ", ")
}

func (m *Method) ReturnArglist() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = p.TypeString()
	}
	if len(m.Returns) > 1 {
		return fmt.Sprintf("(%s)", strings.Join(params, ", "))
	}
	return strings.Join(params, ", ")
}

type Param struct {
	Name     string
	Type     string
	Variadic bool
}

func (p Param) String() string {
	return fmt.Sprintf("%s %s", p.Name, p.TypeString())
}

func (p Param) CallName() string {
	if p.Variadic {
		return p.Name + "..."
	}
	return p.Name
}

func (p Param) TypeString() string {
	if p.Variadic {
		return "..." + p.Type[2:]
	}
	return p.Type
}

type qualifier struct {
	pkgPath string
	imports map[string]struct{}
}

func newQualifier(pkgPath string) *qualifier {
	return &qualifier{
		pkgPath: pkgPath,
		imports: make(map[string]struct{}),
	}
}

func (q *qualifier) Func(pkg *types.Package) string {
	if q.pkgPath != "" && q.pkgPath == pkg.Path() {
		return ""
	}

	path := pkg.Path()
	if pkg.Path() == "." {
		wd, err := os.Getwd()
		if err == nil {
			path = stripGopath(wd)
		}
	}
	q.imports[path] = struct{}{}

	return pkg.Name()
}

func (q *qualifier) Imports() (imports []string) {
	for pkgToImport := range q.imports {
		imports = append(imports, stripVendorPath(pkgToImport))
	}
	return
}

type Result struct {
	SrcPkgPrefix string
	PkgName      string
	Imports      []string
	Interface    *Interface
}

type Interface struct {
	Name    string
	Methods []*Method
}

func ReflectInterface(srcDir, pkgName, objName string) (*Result, error) {
	srcPkgType, pkgPath, pkgName := GetPkgInfo(srcDir, pkgName)
	obj := srcPkgType.Scope().Lookup(objName)
	if obj == nil {
		return nil, fmt.Errorf("cannot find interface %s", objName)
	}

	if !types.IsInterface(obj.Type()) {
		return nil, fmt.Errorf("%s (%s) not an interface", objName, obj.Type().String())
	}

	qualifier := newQualifier(pkgPath)
	var methods []*Method

	ifObj := obj.Type().Underlying().(*types.Interface).Complete()
	for i := 0; i < ifObj.NumMethods(); i++ {
		meth := ifObj.Method(i)
		sig := meth.Type().(*types.Signature)

		method := &Method{
			Name:    meth.Name(),
			Params:  extractArgs(qualifier.Func, sig, sig.Params(), "in%d"),
			Returns: extractArgs(qualifier.Func, sig, sig.Results(), "out%d"),
		}
		methods = append(methods, method)
	}

	srcPkgPrefix := ""
	imports := qualifier.Imports()
	if srcPkgType.Name() != pkgName {
		srcPkgPrefix = srcPkgType.Name() + "."
		imports = append(imports, stripVendorPath(srcPkgType.Path()))
	}

	return &Result{
		SrcPkgPrefix: srcPkgPrefix,
		PkgName:      pkgName,
		Imports:      imports,
		Interface: &Interface{
			Name:    objName,
			Methods: methods,
		},
	}, nil
}

func GetPkgInfo(srcDir, pkgName string) (*types.Package, string, string) {
	srcPkg, err := pkgInfoFromPath(srcDir, packages.NeedName|packages.NeedTypes|packages.NeedTypesInfo)
	if err != nil {
		panic(fmt.Errorf("couldn't load source package: %s", err))
	}

	pkgPath, err := findPkgPath(pkgName, srcPkg)
	if err != nil {
		panic(fmt.Errorf("couldn't load package: %s", err))
	}

	if pkgName == "" {
		pkgName = srcPkg.Name
	}

	return srcPkg.Types, pkgPath, pkgName
}

func extractArgs(qf types.Qualifier, sig *types.Signature, list *types.Tuple, nameFormat string) []*Param {
	var params []*Param
	listLen := list.Len()
	for ii := 0; ii < listLen; ii++ {
		p := list.At(ii)
		name := p.Name()
		if name == "" {
			name = fmt.Sprintf(nameFormat, ii+1)
		}
		typename := types.TypeString(p.Type(), qf)
		// check for final variadic argument
		variadic := sig.Variadic() && ii == listLen-1 && typename[0:2] == "[]"
		param := &Param{
			Name:     name,
			Type:     typename,
			Variadic: variadic,
		}
		params = append(params, param)
	}
	return params
}

// stripVendorPath strips the vendor dir prefix from a package path.
// For example we might encounter an absolute path like
// github.com/foo/bar/vendor/github.com/pkg/errors which is resolved
// to github.com/pkg/errors.
func stripVendorPath(p string) string {
	parts := strings.Split(p, "/vendor/")
	if len(parts) == 1 {
		return p
	}
	return strings.TrimLeft(path.Join(parts[1:]...), "/")
}

// stripGopath takes the directory to a package and removes the
// $GOPATH/src path to get the canonical package name.
func stripGopath(p string) string {
	for _, srcDir := range build.Default.SrcDirs() {
		rel, err := filepath.Rel(srcDir, p)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		return filepath.ToSlash(rel)
	}
	return p
}

func findPkgPath(pkgInputVal string, srcPkg *packages.Package) (string, error) {
	if pkgInputVal == "" {
		return srcPkg.PkgPath, nil
	}
	if pkgInDir(".", pkgInputVal) {
		return ".", nil
	}
	if pkgInDir(srcPkg.PkgPath, pkgInputVal) {
		return srcPkg.PkgPath, nil
	}
	subdirectoryPath := filepath.Join(srcPkg.PkgPath, pkgInputVal)
	if pkgInDir(subdirectoryPath, pkgInputVal) {
		return subdirectoryPath, nil
	}
	return "", nil
}

func pkgInDir(pkgName, dir string) bool {
	currentPkg, err := pkgInfoFromPath(dir, packages.NeedName)
	if err != nil {
		return false
	}
	return currentPkg.Name == pkgName || currentPkg.Name+"_test" == pkgName
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
