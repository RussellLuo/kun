package profilesvc

import (
	"net/http"
)

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func errorToResponse(err error) (int, interface{}) {
	return codeFrom(err), map[string]string{
		"error": err.Error(),
	}
}
