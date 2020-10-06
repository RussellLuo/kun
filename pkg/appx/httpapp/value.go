package httpapp

import (
	"fmt"

	"github.com/RussellLuo/appx"
	"github.com/go-chi/chi"
)

// Value holds attributes of an HTTP application in Go kit.
type Value struct {
	Service interface{} // The Go kit service.
	Router  chi.Router  // The HTTP router.
}

func MustGetService(value interface{}) interface{} {
	return value.(*Value).Service
}

func GetRouter(value interface{}) (chi.Router, error) {
	val, ok := value.(*Value)
	if !ok {
		return nil, fmt.Errorf("value %#v cannot be converted to *httpapp.Value", value)
	}

	if val == nil || val.Router == nil {
		return nil, fmt.Errorf("value %#v is not routable", val)
	}

	return val.Router, nil
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

// RequiredServiceGetter is a helper that makes it easy to get the service
// from a required application, which is bound to a context.
type RequiredServiceGetter struct {
	ctx appx.Context
}

func R(ctx appx.Context) *RequiredServiceGetter {
	return &RequiredServiceGetter{ctx: ctx}
}

func (g *RequiredServiceGetter) MustGet(name string) interface{} {
	return MustGetService(g.ctx.Required[name].Value)
}
