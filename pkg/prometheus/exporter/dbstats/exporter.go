package dbstats

import (
	"database/sql"
	"errors"
	"time"

	"github.com/RussellLuo/kok/pkg/prometheus/metric"
	"github.com/RussellLuo/kok/pkg/tickdoer"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrBadDB    = errors.New("db is nil")
	ErrDBExists = errors.New("db already exists")
)

type Opts struct {
	Namespace string
	Subsystem string

	UpdateInterval time.Duration // How often we fetch and update the DBStats.
}

// Metrics for sql.DBStats.
type Metrics struct {
	MaxOpenConnections *prometheus.GaugeVec // Maximum number of open connections to the database.

	// Pool Status
	OpenConnections *prometheus.GaugeVec // The number of established connections both in use and idle.
	InUse           *prometheus.GaugeVec // The number of connections currently in use.
	Idle            *prometheus.GaugeVec // The number of idle connections.

	// Counters
	WaitCount         *metric.SettableCounterVec // The total number of connections waited for.
	WaitDuration      *metric.SettableCounterVec // The total time (in seconds) blocked waiting for a new connection.
	MaxIdleClosed     *metric.SettableCounterVec // The total number of connections closed due to SetMaxIdleConns.
	MaxLifetimeClosed *metric.SettableCounterVec // The total number of connections closed due to SetConnMaxLifetime.
}

// Exporter exports sql.DBStats as Prometheus metrics.
type Exporter struct {
	opts     *Opts
	metrics  *Metrics
	dbLabels map[*sql.DB]prometheus.Labels

	doer *tickdoer.TickDoer
}

func NewExporter(opts *Opts, labelNames []string) *Exporter {
	return &Exporter{
		opts: opts,
		metrics: &Metrics{
			MaxOpenConnections: metric.NewGaugeVecFrom(prometheus.GaugeOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "max_open_connections",
				Help:      "Maximum number of open connections to the database.",
			}, labelNames),
			OpenConnections: metric.NewGaugeVecFrom(prometheus.GaugeOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "open_connections",
				Help:      "The number of established connections both in use and idle.",
			}, labelNames),
			InUse: metric.NewGaugeVecFrom(prometheus.GaugeOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "in_use",
				Help:      "The number of connections currently in use.",
			}, labelNames),
			Idle: metric.NewGaugeVecFrom(prometheus.GaugeOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "idle",
				Help:      "The number of idle connections.",
			}, labelNames),
			WaitCount: metric.NewSettableCounterVecFrom(metric.SettableCounterOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "wait_count",
				Help:      "The total number of connections waited for.",
			}, labelNames),
			WaitDuration: metric.NewSettableCounterVecFrom(metric.SettableCounterOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "wait_duration",
				Help:      "The total time (in seconds) blocked waiting for a new connection.",
			}, labelNames),
			MaxIdleClosed: metric.NewSettableCounterVecFrom(metric.SettableCounterOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "max_idle_closed",
				Help:      "The total number of connections closed due to SetMaxIdleConns.",
			}, labelNames),
			MaxLifetimeClosed: metric.NewSettableCounterVecFrom(metric.SettableCounterOpts{
				Namespace: opts.Namespace,
				Subsystem: opts.Subsystem,
				Name:      "max_lifetime_closed",
				Help:      "The total number of connections closed due to SetConnMaxLifetime.",
			}, labelNames),
		},
		dbLabels: make(map[*sql.DB]prometheus.Labels),
	}
}

func (e *Exporter) MustBind(db *sql.DB, labelValues ...string) {
	if db == nil {
		panic(ErrBadDB)
	}
	if _, ok := e.dbLabels[db]; ok {
		panic(ErrDBExists)
	}
	e.dbLabels[db] = metric.MakeLabels(labelValues...)
}

func (e *Exporter) update() {
	for db, labels := range e.dbLabels {
		stats := db.Stats()

		e.metrics.MaxOpenConnections.With(labels).Set(float64(stats.MaxOpenConnections))
		e.metrics.OpenConnections.With(labels).Set(float64(stats.OpenConnections))
		e.metrics.InUse.With(labels).Set(float64(stats.InUse))
		e.metrics.Idle.With(labels).Set(float64(stats.Idle))
		e.metrics.WaitCount.With(labels).Set(float64(stats.WaitCount))
		e.metrics.WaitDuration.With(labels).Set(stats.WaitDuration.Seconds())
		e.metrics.MaxIdleClosed.With(labels).Set(float64(stats.MaxIdleClosed))
		e.metrics.MaxLifetimeClosed.With(labels).Set(float64(stats.MaxLifetimeClosed))
	}
}

func (e *Exporter) Start() {
	if e.doer == nil {
		e.doer = tickdoer.TickFunc(e.opts.UpdateInterval, e.update)
	}
}

func (e *Exporter) Stop() {
	e.doer.Stop()
}
