package werror

import (
	"fmt"
)

type Error struct {
	Err error // the underlying error

	// The error message of the underlying error, or an empty string if
	// the underlying error is nil.
	Code    string
	Message string
}

// Wrap wraps err with a new error, whose error message is inherited from msgErr.
func Wrap(err, msgErr error) *Error {
	return wrap(err, msgErr.Error())
}

// Wrapf wraps err with a new error, whose error message is calculated by formatting.
func Wrapf(err error, format string, a ...interface{}) *Error {
	return wrap(err, fmt.Sprintf(format, a...))
}

func wrap(err error, msg string) *Error {
	code := ""
	if err != nil {
		code = err.Error()
	}

	return &Error{
		Err:     err,
		Code:    code,
		Message: msg,
	}
}

// DEPRECATED
func (e *Error) SetError(err error) *Error {
	e.Message = err.Error()
	return e
}

// DEPRECATED
func (e *Error) SetErrorf(format string, a ...interface{}) *Error {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Message }

// Unwrap follows the Unwrap convention introduced in Go 1.13,
// See https://blog.golang.org/go1.13-errors
func (e *Error) Unwrap() error { return e.Err }
