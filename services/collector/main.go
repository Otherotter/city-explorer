package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Otherotter/city-explorer/shared/observability"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Otherotter/city-explorer/services/collector/internal/config"
	"github.com/Otherotter/city-explorer/services/collector/internal/db"
	"github.com/Otherotter/city-explorer/services/collector/internal/fetchers/overpass"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

// collectorHandler holds dependencies for the collect endpoint.
// This is dependency injection — metrics and db passed in,
// not accessed as globals.
type collectorHandler struct {
	db      *sql.DB
	metrics *observability.CollectorMetrics
}

var database *sql.DB

func main() {
	ctx := context.Background()

	slog.SetDefault(observability.NewLogger("collector"))
	slog.Info("collector service starting...")
	// Initialize tracer
	shutdown, traceerr := observability.InitTracer(ctx, "city-explorer-collector")

	if traceerr != nil {
		slog.Error("failed to init tracer", "error", traceerr)
		os.Exit(1)
	}
	defer shutdown(ctx)
	slog.Info("tracer initialized")

	// Try local dev path first, fall back silently
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}
	slog.Info("environment loaded")

	var err error
	database, err = db.Connect()
	if err != nil {
		// log.Fatalf("[collector] could not connect to database: %v", err)
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	// log.Println("[collector] connected to database")
	slog.Info("database connected")

	// Initialize metrics once at startup
	// promauto registers them with the default registry
	// promhttp.Handler() serves them at /metrics
	metrics := observability.NewCollectorMetrics()
	h := &collectorHandler{
		db:      database,
		metrics: metrics,
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(observability.RequestLogger)
	r.Get("/health", healthHandler)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/collect/{city}/{category}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "collect",
				otelhttp.WithMessageEvents(
					otelhttp.ReadEvents,
					otelhttp.WriteEvents,
				),
			)
		})
		r.Post("/", h.collectHandler)
	})

	slog.Info("collector listening", "port", 8081)
	//log.Println("[collector] starting on port 8081") // Old Logging

	if err := http.ListenAndServe(":8081", r); err != nil {
		// log.Fatalf("[collector] failed: %v", err)
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Service:   "collector",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *collectorHandler) collectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("city-explorer-collector")

	citySlug := chi.URLParam(r, "city")
	categorySlug := chi.URLParam(r, "category")
	start := time.Now()

	slog.Info("collection triggered",
		"city", citySlug,
		"category", categorySlug,
	)

	// Span 1 — city lookup
	ctx, citySpan := tracer.Start(ctx, "db.city_lookup")
	citySpan.SetAttributes(attribute.String("city", citySlug))

	var cityID int
	var bbox string
	err := h.db.QueryRowContext(ctx, `
        SELECT id, COALESCE(bbox, '')
        FROM cities
        WHERE LOWER(name) = LOWER($1)
        OR LOWER(name) = LOWER(REPLACE($1, '-', ' '))
    `, citySlug).Scan(&cityID, &bbox)

	if err == sql.ErrNoRows {
		citySpan.SetStatus(codes.Error, "city not found")
		citySpan.RecordError(err)
		citySpan.End()
		slog.Warn("city not found", "city", citySlug)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "city not found", http.StatusNotFound)
		return
	}
	if err != nil {
		citySpan.SetStatus(codes.Error, "db error")
		citySpan.RecordError(err)
		citySpan.End()
		slog.Error("db error looking up city",
			"city", citySlug, "error", err)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "internal server error",
			http.StatusInternalServerError)
		return
	}
	if bbox == "" {
		citySpan.SetStatus(codes.Error, "no bounding box")
		citySpan.End()
		slog.Warn("city has no bounding box",
			"city", citySlug, "city_id", cityID)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "city has no bounding box configured",
			http.StatusUnprocessableEntity)
		return
	}
	citySpan.SetAttributes(
		attribute.Int("city.id", cityID),
		attribute.String("bbox", bbox),
	)
	citySpan.End()

	// Span 2 — category lookup
	ctx, catSpan := tracer.Start(ctx, "db.category_lookup")
	catSpan.SetAttributes(attribute.String("category", categorySlug))

	var categoryID int
	err = h.db.QueryRowContext(ctx, `
        SELECT id FROM categories WHERE slug = $1
    `, categorySlug).Scan(&categoryID)

	if err == sql.ErrNoRows {
		catSpan.SetStatus(codes.Error, "category not found")
		catSpan.RecordError(err)
		catSpan.End()
		slog.Warn("category not found", "category", categorySlug)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		catSpan.SetStatus(codes.Error, "db error")
		catSpan.RecordError(err)
		catSpan.End()
		slog.Error("db error looking up category",
			"category", categorySlug, "error", err)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "internal server error",
			http.StatusInternalServerError)
		return
	}
	catSpan.SetAttributes(attribute.Int("category.id", categoryID))
	catSpan.End()

	// Span 3 — get amenities
	amenities := config.AmenitiesForCategory(categorySlug)
	if len(amenities) == 0 {
		slog.Warn("no amenities configured",
			"category", categorySlug)
		h.metrics.CollectionErrors.
			WithLabelValues("config", categorySlug).Inc()
		http.Error(w, "no amenities configured",
			http.StatusUnprocessableEntity)
		return
	}

	slog.Info("fetching from overpass",
		"city", citySlug,
		"category", categorySlug,
		"amenity_count", len(amenities),
		"bbox", bbox,
	)

	// Span 4 — overpass fetch
	// FetchPlaces creates overpass.http_request as a child
	ctx, fetchSpan := tracer.Start(ctx, "overpass.fetch_places")
	fetchSpan.SetAttributes(
		attribute.String("city", citySlug),
		attribute.String("category", categorySlug),
		attribute.Int("amenity_count", len(amenities)),
	)

	overpassStart := time.Now()
	places, err := overpass.FetchPlaces(
		ctx, cityID, categoryID, amenities, bbox)
	overpassDuration := time.Since(overpassStart)

	h.metrics.OverpassDuration.
		WithLabelValues(categorySlug).
		Observe(overpassDuration.Seconds())

	if err != nil {
		fetchSpan.SetStatus(codes.Error, "overpass fetch failed")
		fetchSpan.RecordError(err)
		fetchSpan.End()
		slog.Error("overpass fetch failed",
			"city", citySlug,
			"category", categorySlug,
			"error", err,
			"duration_ms", overpassDuration.Milliseconds(),
		)
		h.metrics.OverpassRequests.
			WithLabelValues("error", categorySlug).Inc()
		h.metrics.CollectionErrors.
			WithLabelValues("overpass", categorySlug).Inc()
		http.Error(w, "collection failed",
			http.StatusInternalServerError)
		return
	}

	h.metrics.OverpassRequests.
		WithLabelValues("success", categorySlug).Inc()

	fetchSpan.SetAttributes(
		attribute.Int("places.fetched", len(places)),
		attribute.Int64("duration_ms", overpassDuration.Milliseconds()),
	)
	fetchSpan.End()

	slog.Info("overpass fetch complete",
		"city", citySlug,
		"category", categorySlug,
		"fetched", len(places),
		"duration_ms", overpassDuration.Milliseconds(),
	)

	// Span 5 — database insert
	ctx, insertSpan := tracer.Start(ctx, "db.insert_places")
	insertSpan.SetAttributes(
		attribute.Int("places.to_insert", len(places)),
	)

	inserted, err := db.InsertPlaces(h.db, places)
	if err != nil {
		insertSpan.SetStatus(codes.Error, "insert failed")
		insertSpan.RecordError(err)
		insertSpan.End()
		slog.Error("db insert failed",
			"city", citySlug,
			"category", categorySlug,
			"fetched", len(places),
			"error", err,
		)
		h.metrics.CollectionErrors.
			WithLabelValues("db_insert", categorySlug).Inc()
		http.Error(w, "database write failed",
			http.StatusInternalServerError)
		return
	}

	h.metrics.PlacesInserted.
		WithLabelValues(citySlug, categorySlug).
		Add(float64(inserted))

	insertSpan.SetAttributes(
		attribute.Int("places.inserted", inserted),
	)
	insertSpan.End()

	slog.Info("collection complete",
		"city", citySlug,
		"category", categorySlug,
		"fetched", len(places),
		"inserted", inserted,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"city":     citySlug,
		"category": categorySlug,
		"fetched":  len(places),
		"inserted": inserted,
	})
}
