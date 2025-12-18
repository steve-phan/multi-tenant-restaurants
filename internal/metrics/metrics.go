package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Database metrics
	DBQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// Business metrics
	OrdersCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
		[]string{"restaurant_id", "status"},
	)

	ReservationsCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reservations_created_total",
			Help: "Total number of reservations created",
		},
		[]string{"restaurant_id", "status"},
	)

	MenuItemsViewedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "menu_items_viewed_total",
			Help: "Total number of menu item views",
		},
		[]string{"restaurant_id"},
	)

	// Authentication metrics
	AuthAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"status"},
	)

	ActiveSessions = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Number of active user sessions",
		},
	)

	// Error metrics
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "handler"},
	)

	// S3 metrics
	S3UploadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "s3_uploads_total",
			Help: "Total number of S3 uploads",
		},
		[]string{"status"},
	)

	S3UploadDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "s3_upload_duration_seconds",
			Help:    "S3 upload duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// IncrementHTTPRequest records an HTTP request
func IncrementHTTPRequest(method, path, status string) {
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
}

// RecordDBQuery records a database query
func RecordDBQuery(operation, table string, duration float64) {
	DBQueriesTotal.WithLabelValues(operation, table).Inc()
	DBQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// IncrementOrdersCreated increments the orders counter
func IncrementOrdersCreated(restaurantID, status string) {
	OrdersCreatedTotal.WithLabelValues(restaurantID, status).Inc()
}

// IncrementReservationsCreated increments the reservations counter
func IncrementReservationsCreated(restaurantID, status string) {
	ReservationsCreatedTotal.WithLabelValues(restaurantID, status).Inc()
}

// IncrementMenuItemViewed increments the menu item view counter
func IncrementMenuItemViewed(restaurantID string) {
	MenuItemsViewedTotal.WithLabelValues(restaurantID).Inc()
}

// IncrementAuthAttempt increments the auth attempts counter
func IncrementAuthAttempt(status string) {
	AuthAttemptsTotal.WithLabelValues(status).Inc()
}

// SetActiveSessions sets the active sessions gauge
func SetActiveSessions(count float64) {
	ActiveSessions.Set(count)
}

// IncrementError increments the error counter
func IncrementError(errorType, handler string) {
	ErrorsTotal.WithLabelValues(errorType, handler).Inc()
}

// IncrementS3Upload increments the S3 upload counter
func IncrementS3Upload(status string) {
	S3UploadsTotal.WithLabelValues(status).Inc()
}

// RecordS3UploadDuration records S3 upload duration
func RecordS3UploadDuration(duration float64) {
	S3UploadDuration.Observe(duration)
}
