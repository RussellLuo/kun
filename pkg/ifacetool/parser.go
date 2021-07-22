package ifacetool

import (
	"fmt"
	"go/types"
	"strings"
)

type Import struct {
	Alias string
	Path  string
}

func (i *Import) ImportString() string {
	s := fmt.Sprintf("%q", i.Path)
	if i.Alias != "" {
		s = i.Alias + " " + s
	}
	return s
}

type Param struct {
	Name       string
	TypeString string
	Type       types.Type `json:"-"`
	Variadic   bool
}

// MethodArg is the representation of the parameter in the function
// signature, e.g. 'name a.Type'.
func (p *Param) MethodArg() string {
	if p.Variadic {
		return fmt.Sprintf("%s ...%s", p.Name, p.TypeString[2:])
	}
	return fmt.Sprintf("%s %s", p.Name, p.TypeString)
}

// CallName returns the string representation of the parameter to be
// used for a method call. For a variadic parameter, it will be of the
// format 'foos...'.
func (p *Param) CallName() string {
	if p.Variadic {
		return p.Name + "..."
	}
	return p.Name
}

type Method struct {
	Name    string
	Params  []*Param
	Returns []*Param
}

// ArgList is the string representation of method parameters, e.g.
// 's string, n int, foo bar.Baz'.
func (m *Method) ArgList() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.MethodArg()
	}
	return strings.Join(params, ", ")
}

// CallArgList is the string representation of method call parameters,
// e.g. 's, n, foo'. In case of a last variadic parameter, it will be of
// the format 's, n, foos...'
func (m *Method) CallArgList() string {
	params := make([]string, len(m.Params))
	for i, p := range m.Params {
		params[i] = p.CallName()
	}
	return strings.Join(params, ", ")
}

// ReturnArgTypeList is the string representation of types returned by method,
// e.g. 'bar.Baz', '(string, error)'.
func (m *Method) ReturnArgTypeList() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = p.TypeString
	}
	if len(m.Returns) > 1 {
		return fmt.Sprintf("(%s)", strings.Join(params, ", "))
	}
	return strings.Join(params, ", ")
}

// ReturnArgValueList is the string representation of values returned
// by method, e.g. 'foo', 's, err'.
func (m *Method) ReturnArgValueList() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = p.Name
	}
	return strings.Join(params, ", ")
}

// ReturnArgNamedValueList is the string representation of named return values
// returned by method, e.g. '(baz bar.Baz)', '(s string, err error)'.
func (m *Method) ReturnArgNamedValueList() string {
	params := make([]string, len(m.Returns))
	for i, p := range m.Returns {
		params[i] = fmt.Sprintf("%s %s", p.Name, p.TypeString)
	}
	return fmt.Sprintf("(%s)", strings.Join(params, ", "))
}

type Data struct {
	PkgName         string
	SrcPkgName      string
	SrcPkgQualifier string
	InterfaceName   string
	Imports         []*Import
	Methods         []*Method
}

type Parser interface {
	Parse(string) (*Data, error)
}
