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
	r.Post("/collect/{city}/{category}", h.collectHandler)
	// r.Post("/collect/{city}/{category}", collectHandler)
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
	// citySlug := chi.URLParam(r, "city")
	// categorySlug := chi.URLParam(r, "category")
	citySlug := chi.URLParam(r, "city")
	categorySlug := chi.URLParam(r, "category")
	start := time.Now()

	// Old Logging
	// log.Printf("[collector] triggered: city=%s category=%s", citySlug, categorySlug)
	slog.Info("collection triggered",
		"city_name", citySlug, // note: city_name not city
		"cat", categorySlug, // note: cat not category
	)

	// Step 1 — look up city from DB
	// We need the city ID and bounding box
	var cityID int
	var bbox string
	err := database.QueryRow(`
		SELECT id, COALESCE(bbox, '')
		FROM cities
		WHERE LOWER(name) = LOWER($1)
		OR LOWER(name) = LOWER(REPLACE($1, '-', ' '))
	`, citySlug).Scan(&cityID, &bbox)

	if err == sql.ErrNoRows {
		slog.Warn("city not found", "city", citySlug)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "city not found", http.StatusNotFound)
		return
	}

	if err != nil {
		slog.Error("db error looking up city",
			"city", citySlug,
			"error", err,
		)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if bbox == "" {
		slog.Warn("city has no bounding box",
			"city", citySlug,
			"city_id", cityID,
		)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "city has no bounding box configured", http.StatusUnprocessableEntity)
		return
	}

	// Step 2 — look up category
	var categoryID int
	err = h.db.QueryRow(`
		SELECT id FROM categories WHERE slug = $1
	`, categorySlug).Scan(&categoryID)

	if err == sql.ErrNoRows {
		slog.Warn("category not found", "category", categorySlug)
		h.metrics.CollectionErrors.
			WithLabelValues("db_lookup", categorySlug).Inc()
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		// log.Printf("[collector] db error looking up category: %v", err)
		slog.Error("db error looking up category",
			"category", categorySlug,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Step 3 — get amenities
	amenities := config.AmenitiesForCategory(categorySlug)
	if len(amenities) == 0 {
		slog.Warn("no amenities configured",
			"category", categorySlug,
		)
		h.metrics.CollectionErrors.
			WithLabelValues("config", categorySlug).Inc()
		http.Error(w, "no amenities configured for this category", http.StatusUnprocessableEntity)
		return
	}

	slog.Info("fetching from overpass",
		"city", citySlug,
		"category", categorySlug,
		"amenity_count", len(amenities),
		"bbox", bbox,
	)
	// Step 4 — fetch from Overpass
	// Time this separately from total collection time
	overpassStart := time.Now()
	places, err := overpass.FetchPlaces(r.Context(), cityID, categoryID, amenities, bbox)
	overpassDuration := time.Since(overpassStart)

	// Always record duration even on failure
	h.metrics.OverpassDuration.
		WithLabelValues(categorySlug).
		Observe(overpassDuration.Seconds())

	if err != nil {
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
		http.Error(w, "collection failed", http.StatusInternalServerError)
		return
	}

	h.metrics.OverpassRequests.
		WithLabelValues("success", categorySlug).Inc()

	slog.Info("overpass fetch complete",
		"city", citySlug,
		"category", categorySlug,
		"fetched", len(places),
		"duration_ms", overpassDuration.Milliseconds(),
	)

	// Step 5 — insert into DB
	inserted, err := db.InsertPlaces(h.db, places)
	if err != nil {
		slog.Error("db insert failed",
			"city", citySlug,
			"category", categorySlug,
			"fetched", len(places),
			"error", err,
		)
		h.metrics.CollectionErrors.
			WithLabelValues("db_insert", categorySlug).Inc()
		http.Error(w, "database write failed", http.StatusInternalServerError)
		return
	}

	h.metrics.PlacesInserted.
		WithLabelValues(citySlug, categorySlug).
		Add(float64(inserted))

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
