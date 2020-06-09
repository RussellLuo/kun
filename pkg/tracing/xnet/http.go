package xnet

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

// HTTPToContext returns an http RequestFunc that associates ctx with a Tracer.
func HTTPToContext(newTracer NewFunc, family, title string) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		tracer := newTracer(family, title)
		return NewContext(ctx, tracer)
	}
}
