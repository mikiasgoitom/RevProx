package metricsadapter

import (
	"context"
	"time"

	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusAdapter implements the IMetricsUseCase interface using Prometheus collectors.
type PrometheusAdapter struct {
	hits      prometheus.Counter
	misses    prometheus.Counter
	evictions prometheus.Counter
	latencies *prometheus.HistogramVec
}

// NewPrometheusAdapter creates and registers the Prometheus metrics.
func NewPrometheusAdapter() contract.IMetricsAdapter {
	return &PrometheusAdapter{
		hits: promauto.NewCounter(prometheus.CounterOpts{
			Name: "caching_proxy_cache_hits_total",
			Help: "The total number of cache hits.",
		}),
		misses: promauto.NewCounter(prometheus.CounterOpts{
			Name: "caching_proxy_cache_misses_total",
			Help: "The total number of cache misses.",
		}),
		evictions: promauto.NewCounter(prometheus.CounterOpts{
			Name: "caching_proxy_cache_evictions_total",
			Help: "The total number of items evicted from the cache.",
		}),
		latencies: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "caching_proxy_latency_seconds",
			Help:    "Request latency in seconds, partitioned by type.",
			Buckets: prometheus.DefBuckets, // Default buckets are fine for now
		}, []string{"type"}), // Labels: "total", "upstream", "cache"
	}
}

func (a *PrometheusAdapter) IncHit(ctx context.Context) error {
	a.hits.Inc()
	return nil
}

func (a *PrometheusAdapter) IncMiss(ctx context.Context) error {
	a.misses.Inc()
	return nil
}

func (a *PrometheusAdapter) RecordEviction(ctx context.Context) error {
	a.evictions.Inc()
	return nil
}

func (a *PrometheusAdapter) RecordUpstreamLatency(ctx context.Context, d time.Duration) error {
	a.latencies.WithLabelValues("upstream").Observe(d.Seconds())
	return nil
}

func (a *PrometheusAdapter) RecordCacheLatency(ctx context.Context, d time.Duration) error {
	a.latencies.WithLabelValues("cache").Observe(d.Seconds())
	return nil
}

func (a *PrometheusAdapter) RecordTotalLatency(ctx context.Context, d time.Duration) error {
	a.latencies.WithLabelValues("total").Observe(d.Seconds())
	return nil
}
