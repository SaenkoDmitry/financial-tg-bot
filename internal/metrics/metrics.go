package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	IncomingRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ozon",
			Subsystem: "http",
			Name:      "incoming_requests_total_counter",
		},
		[]string{"model_type", "command", "status"},
	)
	IncomingRequestsHistogramResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ozon",
			Subsystem: "http",
			Name:      "incoming_requests_histogram_response_time_seconds",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.5, 1, 2},
		},
		[]string{"model_type", "command", "status"},
	)
	CacheHitCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ozon",
			Subsystem: "http",
			Name:      "cache_hit_counter",
		},
		[]string{"status"},
	)
	RatesSourceCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ozon",
			Subsystem: "http",
			Name:      "rates_source_counter",
		},
		[]string{"source"},
	)
	RatesAPICallCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "ozon",
			Subsystem: "http",
			Name:      "rates_api_call_counter",
		},
		[]string{"type", "status"},
	)
)

var (
	HitLabel  = "hit"
	MissLabel = "miss"
)

var (
	CacheLabel = "cache"
	DBLabel    = "database"
	APILabel   = "api"
)

var (
	LiveCallTypeLabel       = "live"
	HistoricalCallTypeLabel = "historical"
)
