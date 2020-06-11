package googlecode

import (
	"errors"
	"net/http"

	"github.com/RussellLuo/kok/pkg/werror"
)

func HTTPStatusCode(err error) int {
	if errors.Is(err, ErrInvalidArgument) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrFailedPrecondition) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrOutOfRange) {
		return http.StatusBadRequest
	} else if errors.Is(err, ErrUnauthenticated) {
		return http.StatusUnauthorized
	} else if errors.Is(err, ErrPermissionDenied) {
		return http.StatusForbidden
	} else if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, ErrAborted) {
		return http.StatusConflict
	} else if errors.Is(err, ErrAlreadyExists) {
		return http.StatusConflict
	} else if errors.Is(err, ErrResourceExhausted) {
		return http.StatusTooManyRequests
	} else if errors.Is(err, ErrCancelled) {
		return 499 // has no HTTP mapping
	} else if errors.Is(err, ErrDataLoss) {
		return http.StatusInternalServerError
	} else if errors.Is(err, ErrUnknown) {
		return http.StatusInternalServerError
	} else if errors.Is(err, ErrInternal) {
		return http.StatusInternalServerError
	} else if errors.Is(err, ErrNotImplemented) {
		return http.StatusNotImplemented
	} else if errors.Is(err, ErrUnavailable) {
		return http.StatusServiceUnavailable
	} else if errors.Is(err, ErrDeadlineExceeded) {
		return http.StatusGatewayTimeout
	} else {
		return http.StatusInternalServerError
	}
}

func HTTPResponse(err error) (int, interface{}) {
	var e *werror.Error
	var code, message string

	if errors.As(err, &e) {
		code, message = e.Code, e.Message
	} else {
		code, message = ErrUnknown.Error(), err.Error()
	}

	return HTTPStatusCode(err), map[string]map[string]string{
		"error": {
			"code":    code,
			"message": message,
		},
	}
}
