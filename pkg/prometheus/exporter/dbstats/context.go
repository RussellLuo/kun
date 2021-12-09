package dbstats

import (
	"context"
)

type contextKeyE string

var contextKey = contextKeyE("github.com/RussellLuo/kun/pkg/prometheus/exporter/dbstats.Exporter")

// NewContext returns a copy of the parent context
// and associates it with an Exporter.
func NewContext(ctx context.Context, ex *Exporter) context.Context {
	return context.WithValue(ctx, contextKey, ex)
}

// FromContext returns the Exporter bound to the context, if any.
func FromContext(ctx context.Context) (ex *Exporter, ok bool) {
	ex, ok = ctx.Value(contextKey).(*Exporter)
	return
}
