package httpoption

import (
	"github.com/RussellLuo/kun/pkg/oas2"
)

type Options struct {
	requestValidators map[string]Validator
	responseSchema    oas2.Schema
}

func NewOptions(opts ...Option) *Options {
	options := &Options{
		requestValidators: make(map[string]Validator),
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

func (o *Options) RequestValidator(name string) Validator {
	if v, ok := o.requestValidators[name]; ok {
		return v
	}
	return nilValidator
}

func (o *Options) ResponseSchema() oas2.Schema {
	if o.responseSchema != nil {
		return o.responseSchema
	}
	return defaultSchema
}

// Option sets an optional parameter for Options.
type Option func(*Options)

// RequestValidators sets the request validators for Options.
func RequestValidators(validators ...NamedValidator) Option {
	return func(o *Options) {
		for _, v := range validators {
			o.requestValidators[v.Name] = v.Validator
		}
	}
}

// ResponseSchema sets the response schema for Options.
func ResponseSchema(schema oas2.Schema) Option {
	return func(o *Options) {
		o.responseSchema = schema
	}
}

// NamedValidator holds a validator and its corresponding operation name, to
// which the request belongs.
type NamedValidator struct {
	Name      string
	Validator Validator
}

// Op is a shortcut for creating an instance of NamedValidator.
func Op(name string, validator Validator) NamedValidator {
	return NamedValidator{
		Name:      name,
		Validator: validator,
	}
}

var (
	// nilValidator is a validator that does no validation.
	nilValidator = FuncValidator(func(value interface{}) error {
		return nil
	})

	defaultSchema = &oas2.ResponseSchema{}
)
