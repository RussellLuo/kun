package profilesvc

import (
	"errors"
	"net/http"
)

func err2code(err error) int {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, ErrAlreadyExists) || errors.Is(err, ErrInconsistentIDs) {
		return http.StatusBadRequest
	} else {
		return http.StatusInternalServerError
	}
}
