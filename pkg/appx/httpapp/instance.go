package httpapp

import (
	"fmt"

	"github.com/RussellLuo/appx"
	"github.com/go-chi/chi"
)

func MountOn(parent, pattern string) appx.Middleware {
	return func(next appx.Instance) appx.Instance {
		return middleware{
			Standard: appx.Standardize(next),
			parent:   parent,
			pattern:  pattern,
		}
	}
}

type middleware struct {
	appx.Standard
	parent  string
	pattern string
}

func (m middleware) Init(ctx appx.Context) error {
	if err := m.Standard.Init(ctx); err != nil {
		return err
	}

	parent, err := getChiRouter(ctx.MustLoad(m.parent))
	if err != nil {
		return err
	}

	r, err := getChiRouter(m.Standard.Instance())
	if err != nil {
		return err
	}

	MountRouter(parent, m.pattern, r)
	return nil
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

func getChiRouter(instance interface{}) (chi.Router, error) {
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
