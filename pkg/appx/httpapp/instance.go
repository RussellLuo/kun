package httpapp

import (
	"fmt"

	"github.com/go-chi/chi"
)

// Value holds attributes of an HTTP application in Go kit.
type Value struct {
	Service interface{} // The Go kit service.
	Router  chi.Router  // The HTTP router.
}

func MountRouter(parent chi.Router, pattern string, r chi.Router) {
	if pattern == "" {
		extendRouter(parent, r)
	} else {
		parent.Mount(pattern, r)
	}
}

func extendRouter(parent chi.Router, r chi.Router) {
	for _, route := range r.Routes() {
		for method, handler := range route.Handlers {
			parent.Method(method, route.Pattern, handler)
		}
	}
}

type ChiRouter interface {
	Router() chi.Router
}

func GetChiRouter(instance interface{}) (chi.Router, error) {
	r, ok := instance.(ChiRouter)
	if !ok {
		return nil, fmt.Errorf("instance %#v does not implement httpapp.ChiRouter", instance)
	}

	result := r.Router()
	if result == nil {
		return nil, fmt.Errorf("method Router() of instance %#v returns nil", instance)
	}

	return result, nil
}
