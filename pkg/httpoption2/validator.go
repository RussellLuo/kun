package httpoption

import (
	"github.com/RussellLuo/kun/pkg/werror"
	"github.com/RussellLuo/kun/pkg/werror/gcode"
	"github.com/RussellLuo/validating/v2"
)

type Validator interface {
	Validate(value interface{}) error
}

// FuncValidator is an adapter to allow the use of ordinary functions as
// validators. If f is a function with the appropriate signature,
// Func(f) is a Validator that calls f.
type FuncValidator func(value interface{}) error

// Validate calls fv(value).
func (fv FuncValidator) Validate(value interface{}) error {
	return fv(value)
}

func Validate(schema validating.Schema) error {
	errs := validating.Validate(schema)
	if len(errs) == 0 {
		return nil
	}
	return werror.Wrap(gcode.ErrInvalidArgument, errs)
}
