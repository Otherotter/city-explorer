package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CollectorMetrics holds all Prometheus metrics
// for the Collector service.
// Initialize once with NewCollectorMetrics()
// and pass it into handlers via dependency injection.
type CollectorMetrics struct {
	// OverpassRequests counts every call to Overpass
	// Labels: status (success|error), category
	OverpassRequests *prometheus.CounterVec

	// OverpassDuration tracks how long Overpass calls take
	// Labels: category
	OverpassDuration *prometheus.HistogramVec

	// PlacesInserted counts places written to the DB
	// Labels: city, category
	PlacesInserted *prometheus.CounterVec

	// CollectionErrors counts failures by stage
	// Labels: stage (overpass|db_lookup|db_insert|config), category
	CollectionErrors *prometheus.CounterVec
}

// NewCollectorMetrics registers and returns all
// Collector Prometheus metrics.
// Call once in main().
func NewCollectorMetrics() *CollectorMetrics {
	return &CollectorMetrics{
		OverpassRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cityexplorer_collector_overpass_requests_total",
				Help: "Total Overpass API requests by status and category",
			},
			[]string{"status", "category"},
		),

		OverpassDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cityexplorer_collector_overpass_duration_seconds",
				Help:    "Overpass API request duration in seconds",
				Buckets: []float64{1, 5, 10, 15, 20, 30, 60},
			},
			[]string{"category"},
		),

		PlacesInserted: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cityexplorer_collector_places_inserted_total",
				Help: "Total places inserted into the database",
			},
			[]string{"city", "category"},
		),

		CollectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cityexplorer_collector_errors_total",
				Help: "Total collection errors by stage and category",
			},
			[]string{"stage", "category"},
		),
	}
}
