package oas2

import (
	"reflect"
)

type GoStruct struct {
	Name       string
	Properties []Property
}

// GoTypeRegistry is a mapping from "fully qualified type name" to "local type name".
type GoTypeRegistry struct {
	registry map[string]*GoStruct
}

func NewGoTypeRegistry() *GoTypeRegistry {
	return &GoTypeRegistry{registry: make(map[string]*GoStruct)}
}

// Register tries to register the given typ and returns whether
// typ is successfully registered.
func (r *GoTypeRegistry) Register(typ reflect.Type, name string) bool {
	fullName := r.FullName(typ)
	if fullName == "." {
		// Do not register anonymous struct types.
		return true
	}

	//fmt.Printf("typ: %s, fullName: %s\n", typ.String(), fullName)
	if _, ok := r.registry[fullName]; ok {
		return false
	}

	r.registry[fullName] = &GoStruct{Name: name}
	return true
}

func (r *GoTypeRegistry) Properties(typ reflect.Type) []Property {
	fullName := r.FullName(typ)
	if gs, ok := r.registry[fullName]; ok {
		return gs.Properties
	}
	return nil
}

func (r *GoTypeRegistry) SetProperties(typ reflect.Type, properties []Property) {
	fullName := r.FullName(typ)
	if gs, ok := r.registry[fullName]; ok {
		gs.Properties = properties
	}
}

func (r *GoTypeRegistry) FullName(typ reflect.Type) string {
	return typ.PkgPath() + "." + typ.Name()
}

func (r *GoTypeRegistry) Name(typ reflect.Type) string {
	fullName := r.FullName(typ)
	if gs, ok := r.registry[fullName]; ok {
		return gs.Name
	}
	return ""
}
