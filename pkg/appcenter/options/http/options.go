package http

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"

	"github.com/RussellLuo/kok/pkg/appcenter"
)

type Options struct {
	Pattern string
	Router  chi.Router
}

func extendRouter(r chi.Router, appRouter chi.Router) {
	for _, route := range appRouter.Routes() {
		for method, handler := range route.Handlers {
			r.Method(method, route.Pattern, handler)
		}
	}
}

func mountRouter(r chi.Router, options interface{}) error {
	opts, ok := options.(*Options)
	if !ok {
		return fmt.Errorf("%v cannot be converted to *Options", options)
	}

	if opts == nil || opts.Router == nil {
		// The corresponding application is not routable.
		return nil
	}

	if opts.Pattern == "" {
		extendRouter(r, opts.Router)
	} else {
		r.Mount(opts.Pattern, opts.Router)
	}

	return nil
}

func MakeMountFunc(r chi.Router) appcenter.MountFunc {
	return func(ctx context.Context, subApps []*appcenter.App) error {
		for _, subApp := range subApps {
			if subApp.Options != nil {
				if err := mountRouter(r, subApp.Options); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
