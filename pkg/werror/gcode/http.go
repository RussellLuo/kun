package gcode

import (
	"errors"
	"net/http"

	"github.com/RussellLuo/kok/pkg/werror"
)

func HTTPStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, ErrFailedPrecondition):
		return http.StatusBadRequest
	case errors.Is(err, ErrOutOfRange):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthenticated):
		return http.StatusUnauthorized
	case errors.Is(err, ErrPermissionDenied):
		return http.StatusForbidden
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrAborted):
		return http.StatusConflict
	case errors.Is(err, ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, ErrResourceExhausted):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrCancelled):
		return 499 // has no corresponding constant
	case errors.Is(err, ErrDataLoss):
		return http.StatusInternalServerError
	case errors.Is(err, ErrUnknown):
		return http.StatusInternalServerError
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, ErrUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrDeadlineExceeded):
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

func ToCodeMessage(err error) (string, string) {
	var e *werror.Error
	if errors.As(err, &e) {
		return e.Code, e.Message
	}
	return ErrUnknown.Error(), err.Error()
}

func FromCodeMessage(code, message string) error {
	codeErr := werror.Wrapf(nil, code)
	return werror.Wrapf(codeErr, message)
}
