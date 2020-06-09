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

// nilTracer is a fake tracer that traces nothing.
type nilTracer struct{}

// NewNilTracer creates a fake tracer.
func NewNilTracer() Tracer {
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
	return NewNilTracer()
}

// Authorizer determines whether a specific request is permitted to load the
// /debug/requests or /debug/events pages.
type Authorizer func(req *http.Request) (any, sensitive bool)

var (
	AllowLocal = trace.AuthRequest
	AllowAny   = func(req *http.Request) (any, sensitive bool) { return true, true }
)

// Traces responds with traces from the program.
// The package initialization registers it in http.DefaultServeMux
// at /debug/requests.
//
// It performs authorization by running auth.
func Traces(w http.ResponseWriter, req *http.Request, auth Authorizer) {
	any, sensitive := auth(req)
	if !any {
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	trace.Render(w, req, sensitive)
}
