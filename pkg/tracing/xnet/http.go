package xnet

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

// HTTPToContext returns an http RequestFunc that associates ctx with a Tracer.
func HTTPToContext(tracer Tracer) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		return NewContext(ctx, tracer)
	}
}
