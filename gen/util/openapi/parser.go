package openapi

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
)

type annotation struct {
	ArgName     string
	In          string
	Name        string
	Type        string
	Required    bool
	Description string
}

func newParamAnnotation(text, prevParamName string) (*annotation, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty @kok(param)")
	}

	split := strings.SplitN(text, "<", 2)
	name := strings.TrimSpace(split[0])

	if name == "" {
		if prevParamName == "" {
			return nil, fmt.Errorf("found no argument name in: %s", text)
		}
		name = prevParamName
	}

	a := &annotation{
		ArgName: name,
		Name:    name,
	}

	if len(split) == 1 {
		// No value definition after the argument name.
		return a, nil
	}

	value := strings.TrimSpace(split[1])
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, ":") {
			panic(fmt.Errorf("invalid tag part: %s", part))
		}

		split := strings.SplitN(part, ":", 2)
		k, v := split[0], split[1]

		switch k {
		case "in":
			a.In = v

			if err := errorIn(a.In); err != nil {
				return nil, err
			}
		case "name":
			a.Name = v
		case "required":
			a.Required = v == "true"
		case "descr":
			a.Description = v
		default:
			return nil, fmt.Errorf("invalid tag part: %s", part)
		}
	}

	if a.In == InPath {
		// Parameters located in path must be required.
		a.Required = true
	}

	if a.In == InRequest && a.Name != "RemoteAddr" {
		return nil, fmt.Errorf("argument %q tries to extract value from `request.%s`, but only `request.RemoteAddr` is available", a.ArgName, a.Name)
	}

	return a, nil
}

func (a *annotation) Copy() *annotation {
	anno := *a
	return &anno
}

func (a *annotation) SetType(typ string) *annotation {
	a.Type = typ
	return a
}

type Parser struct {
	methodName string
	params     map[string]*Param

	prevArgName string
}

func (p *Parser) Parse(text string) ([]*annotation, error) {
	a, err := newParamAnnotation(text, p.prevArgName)
	if err != nil {
		return nil, err
	}

	param, ok := p.params[a.ArgName]
	if !ok {
		return nil, fmt.Errorf("no param `%s` declared in the method %s", a.ArgName, p.methodName)
	}

	annotations, err := p.completeAnnotations(param.RawType, a)
	if err != nil {
		return nil, err
	}

	// Record the previous argument name.
	p.prevArgName = a.ArgName

	return annotations, nil
}

// completeAnnotationsPerType builds the complete annotations according to the type
// of the method argument and the original annotation (usually a simplified version)
// the user specified.
func (p *Parser) completeAnnotations(typ types.Type, a *annotation) (annotations []*annotation, err error) {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		annotations = append(annotations, a.Copy().SetType(t.Name()))
	case *types.Slice:
		et, ok := t.Elem().(*types.Basic)
		if !ok {
			return nil, p.error(a.ArgName, t)
		}
		annotations = append(annotations, a.Copy().SetType("[]"+et.Name()))
	case *types.Struct:
		/*if a.In != "" {
			fmt.Printf("WARNING: manually specified `in:%s` is ignored for struct argument `%s` in method %s\n", a.In, a.ArgName, p.methodName)
		}*/

		for i := 0; i < t.NumFields(); i++ {
			var typeName string
			switch ft := t.Field(i).Type().(type) {
			case *types.Basic:
				typeName = ft.Name()
			case *types.Slice:
				et, ok := ft.Elem().(*types.Basic)
				if !ok {
					return nil, p.errorStruct(a.ArgName, t.Field(i).Name(), ft)
				}
				typeName = "[]" + et.Name()
			default:
				return nil, p.errorStruct(a.ArgName, t.Field(i).Name(), ft)
			}

			field := reflect.StructField{
				Tag:  reflect.StructTag(t.Tag(i)),
				Name: t.Field(i).Name(),
			}

			kokField := httpcodec.GetKokField(field)
			if kokField.Omitted {
				continue
			}

			parts := strings.SplitN(kokField.Name, ".", 2)
			in, name := parts[0], parts[1]
			if err := errorIn(in); err != nil {
				return nil, err
			}

			anno := a.Copy().SetType(typeName)
			anno.In = in
			anno.Name = name
			anno.Required = kokField.Required
			annotations = append(annotations, anno)
		}
	//case *types.Pointer:
	//	// Dereference the pointer to parse the element type.
	//	nts, err := p.parseTypes(name, t.Elem())
	//	if err != nil {
	//		return nil, err
	//	}
	//	for n, t := range nts {
	//		nameTypes[n] = t
	//	}
	default:
		return nil, p.error(a.ArgName, t)
	}

	return annotations, nil
}

func (p *Parser) error(name string, typ types.Type) error {
	return fmt.Errorf("parameter cannot be mapped to argument `%s` (of type %v) in method %s", name, typ, p.methodName)
}

func (p *Parser) errorStruct(paramName, fieldName string, fieldType types.Type) error {
	return fmt.Errorf("parameter cannot be mapped to struct field %q (of type %v) from argument `%s` in method %s", fieldName, fieldType, paramName, p.methodName)
}

func errorIn(in string) error {
	if in != InPath && in != InQuery && in != InHeader &&
		/*in != InCookie &&*/ in != InRequest {

		return fmt.Errorf(
			"invalid location value: %s (must be %q, %q, %q or %q)",
			in, InPath, InQuery, InHeader, InRequest,
		)
	}
	return nil
}
