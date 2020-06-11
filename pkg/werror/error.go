package werror

import (
	"fmt"
)

type Error struct {
	Err error // the underlying error

	Code    string
	Message string
}

func Wrap(err error) *Error {
	code := ""
	if err != nil {
		code = err.Error()
	}

	return &Error{
		Err:  err,
		Code: code,
	}
}

func (e *Error) SetError(err error) *Error {
	e.Message = err.Error()
	return e
}

func (e *Error) SetErrorf(format string, a ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Message }

// Unwrap follows the Unwrap convention introduced in Go 1.13,
// See https://blog.golang.org/go1.13-errors
func (e *Error) Unwrap() error { return e.Err }
