package xnet

import (
	"fmt"
	"net/http"

	"golang.org/x/net/trace"
)

type Tracer interface {
	trace.Trace
}

// nilTracer is a fake tracer that traces nothing.
type nilTracer struct{}

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
func Traces(auth Authorizer) http.HandlerFunc {
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
