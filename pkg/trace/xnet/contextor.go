package xnet

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"

	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/trace"
)

// NewContext returns a copy of the parent context
// and associates it with a Trace.
func NewContext(ctx context.Context, tr Tracer) context.Context {
	return trace.NewContext(ctx, tr)
}

// FromContext returns the Trace bound to the context, if any.
func FromContext(ctx context.Context) Tracer {
	if tr, ok := trace.FromContext(ctx); ok {
		return tr
	}
	return nilTracer{}
}

// Contextor is a context manager.
type Contextor struct {
	enabled int32
}

func NewContextor() *Contextor {
	return &Contextor{}
}

// Enable enables the request tracing.
func (c *Contextor) Enable() {
	atomic.StoreInt32(&c.enabled, 1)
}

// Disable disables the request tracing.
func (c *Contextor) Disable() {
	atomic.StoreInt32(&c.enabled, 0)
}

// HTTPToContext returns an http RequestFunc that associates ctx with a Tracer.
func (c *Contextor) HTTPToContext(family, title string) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		if atomic.LoadInt32(&c.enabled) != 1 {
			return ctx
		}
		tracer := NewTracer(family, title)
		return NewContext(ctx, tracer)
	}
}

// HTTPHandler returns an HTTP handler that can be used to enable or disable the tracing.
func HTTPHandler(c *Contextor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if body.Enabled {
			c.Enable()
		} else {
			c.Disable()
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
	}
}
