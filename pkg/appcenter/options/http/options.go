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

func getOptions(app *appcenter.App) (*Options, error) {
	if app.Options == nil {
		return nil, nil
	}

	opts, ok := app.Options.(*Options)
	if !ok {
		return nil, fmt.Errorf("options %v cannot be converted to *Options in app: %v", app.Options, app)
	}

	return opts, nil
}

func extendRouter(r chi.Router, appRouter chi.Router) {
	for _, route := range appRouter.Routes() {
		for method, handler := range route.Handlers {
			r.Method(method, route.Pattern, handler)
		}
	}
}

func mountRouter(r chi.Router, app *appcenter.App) error {
	opts, err := getOptions(app)
	if err != nil {
		return err
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

func MountRouter(ctx context.Context, app *appcenter.App, subApps []*appcenter.App) error {
	opts, err := getOptions(app)
	if err != nil {
		return err
	}

	if opts == nil {
		return fmt.Errorf("nil options in app: %v", app)
	}

	r := opts.Router
	if r == nil {
		return fmt.Errorf("nil router in app: %v", app)
	}

	for _, subApp := range subApps {
		if err := mountRouter(r, subApp); err != nil {
			return err
		}
	}
	return nil
}
