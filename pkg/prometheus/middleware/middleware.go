package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type InstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func (m *InstrumentingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func(begin time.Time) {
			rctx := chi.RouteContext(r.Context())
			if rctx != nil {
				// We are actually using chi.

				path := rctx.RoutePattern()
				status := strconv.Itoa(ww.Status())
				m.requestCount.With("code", status, "method", r.Method, "path", path).Add(1)
				m.requestLatency.With("code", status, "method", r.Method, "path", path).Observe(time.Since(begin).Seconds())
			}
		}(time.Now())

		next.ServeHTTP(ww, r)
	})
}

func NewInstrumentingMiddleware(name string) *InstrumentingMiddleware {
	fieldKeys := []string{"code", "method", "path"}
	return &InstrumentingMiddleware{
		requestCount: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: name,
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		requestLatency: kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: name,
			Name:      "request_latency",
			Help:      "Total duration of requests in seconds.",
		}, fieldKeys),
	}
}
