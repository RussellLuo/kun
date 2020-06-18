package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// SettableCounterOpts is an alias for prometheus.GaugeOpts.
type SettableCounterOpts prometheus.GaugeOpts

// SettableCounterVec is a Counter metric that can be updated by using `Set()`.
//
// Some folks suggest implementing custom collector, see
// https://www.robustperception.io/setting-a-prometheus-counter
type SettableCounterVec struct {
	*prometheus.GaugeVec
}

// NewSettableCounterVec constructs and returns a SettableGaugeVec.
func NewSettableCounterVec(opts SettableCounterOpts, labelNames []string) *SettableCounterVec {
	gv := prometheus.NewGaugeVec(prometheus.GaugeOpts(opts), labelNames)
	return &SettableCounterVec{GaugeVec: gv}
}

// NewSettableCounterVecFrom constructs, registers and returns a SettableGaugeVec.
func NewSettableCounterVecFrom(opts SettableCounterOpts, labelNames []string) *SettableCounterVec {
	scv := NewSettableCounterVec(opts, labelNames)
	prometheus.MustRegister(scv)
	return scv
}

// NewGaugeFrom constructs, registers and returns a GaugeVec.
func NewGaugeVecFrom(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return gv
}

func MakeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
