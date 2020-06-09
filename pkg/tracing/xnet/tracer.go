package xnet

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/net/trace"
)

type Tracer interface {
	trace.Trace
}

type NewFunc func(family, title string) Tracer

// nilTracer is a fake tracer that traces nothing.
type nilTracer struct{}

// NewNil creates a fake tracer.
func NewNil(family, title string) Tracer {
	return nilTracer{}
}

func (tr nilTracer) LazyLog(x fmt.Stringer, sensitive bool)     {}
func (tr nilTracer) LazyPrintf(format string, a ...interface{}) {}
func (tr nilTracer) SetError()                                  {}
func (tr nilTracer) SetRecycler(f func(interface{}))            {}
func (tr nilTracer) SetTraceInfo(traceID, spanID uint64)        {}
func (tr nilTracer) SetMaxEvents(m int)                         {}
func (tr nilTracer) Finish()                                    {}

// NewTracer creates a real tracer.
func NewTracer(family, title string) Tracer {
	return trace.New(family, title)
}

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
	return NewNil("", "")
}

// Authorizer determines whether a specific request is permitted to load the
// /debug/requests or /debug/events pages.
type Authorizer func(r *http.Request) (any, sensitive bool)

var (
	AllowLocal = trace.AuthRequest
	AllowAny   = func(r *http.Request) (any, sensitive bool) { return true, true }
)

// Traces returns an HTTP handler, which will respond with traces from the program.
//
// The handler performs authorization by running auth.
func Traces(auth Authorizer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		any, sensitive := auth(r)
		if !any {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.Render(w, r, sensitive)
	}
}
