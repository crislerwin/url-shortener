// Package metrics provides Prometheus metrics for the URL shortener service.
//
// This package defines and registers all metrics collectors for monitoring
// URL shortening operations, redirects, and HTTP requests.
//
// Example usage:
//
//	// Initialize metrics
//	metrics.Init()
//
//	// Record operations
//	metrics.URLsShortenedTotal.Inc()
//	metrics.URLRedirectsTotal.WithLabelValues("success").Inc()
//	metrics.HTTPRequestDuration.WithLabelValues("/shorten", "POST", "200").Observe(0.5)
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// URLsShortenedTotal tracks the total number of URLs shortened
	URLsShortenedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_shortener_urls_shortened_total",
		Help: "The total number of URLs that have been shortened",
	})

	// URLRedirectsTotal tracks the total number of URL redirects
	URLRedirectsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_redirects_total",
			Help: "The total number of URL redirects",
		},
		[]string{"status"}, // status: success, not_found, error
	)

	// URLsStoredTotal tracks the total number of URLs currently stored
	URLsStoredGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "url_shortener_urls_stored",
		Help: "The current number of URLs stored in the database",
	})

	// HTTPRequestDuration tracks HTTP request duration in seconds
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "url_shortener_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	// HTTPRequestsTotal tracks the total number of HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_http_requests_total",
			Help: "The total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	// DatabaseOperationDuration tracks database operation duration
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "url_shortener_db_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"}, // operation: get, set, delete, exists, count
	)

	// DatabaseErrorsTotal tracks database errors
	DatabaseErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_db_errors_total",
			Help: "The total number of database errors",
		},
		[]string{"operation"},
	)

	// EncryptionOperationsTotal tracks encryption/decryption operations
	EncryptionOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_encryption_operations_total",
			Help: "The total number of encryption/decryption operations",
		},
		[]string{"operation", "status"}, // operation: encrypt, decrypt; status: success, error
	)
)

// Init initializes the metrics package
// This function can be used for any future initialization needs
func Init() {
	// Metrics are automatically registered via promauto
	// This function is provided for future extensibility
}
