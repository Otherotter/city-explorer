package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Otherotter/city-explorer/services/collector/internal/config"
	"github.com/Otherotter/city-explorer/services/collector/internal/db"
	"github.com/Otherotter/city-explorer/services/collector/internal/fetchers/overpass"
	"github.com/Otherotter/city-explorer/shared/observability"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

var database *sql.DB

func main() {
	// ++++
	// LOG <
	slog.SetDefault(observability.NewLogger("collector"))
	slog.Info("collector service starting...")
	// Try local dev path first, fall back silently
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}
	slog.Info("environment loaded")
	// > LOG
	// ++++

	// ++++
	// DB <
	var err error
	database, err = db.Connect()
	if err != nil {
		// log.Fatalf("[collector] could not connect to database: %v", err)
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	// log.Println("[collector] connected to database")
	slog.Info("database connected") // New Logging
	// > DB
	// ++++

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.Post("/collect/{city}/{category}", collectHandler)

	//log.Println("[collector] starting on port 8081") // Old Logging
	slog.Info("starting service",
		"port", 8081, // note: city_name not city
	)

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

func collectHandler(w http.ResponseWriter, r *http.Request) {
	citySlug := chi.URLParam(r, "city")
	categorySlug := chi.URLParam(r, "category")

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
		// log.Printf("[collector] city not found: %s", citySlug) // Old Logging
		slog.Error("collection triggered", "city_name", citySlug)
		http.Error(w, "city not found", http.StatusNotFound)
		return
	}
	if err != nil {
		// log.Printf("[collector] db error looking up city: %v", err)
		slog.Error("db error looking up city", "err", err) // Old Logging
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if bbox == "" {
		// log.Printf("[collector] no bounding box for city: %s", citySlug)
		slog.Error(" no bounding box for city", "city_name", citySlug)
		http.Error(w, "city has no bounding box configured", http.StatusUnprocessableEntity)
		return
	}

	// Step 2 — look up category from DB
	var categoryID int
	err = database.QueryRow(`
		SELECT id FROM categories WHERE slug = $1
	`, categorySlug).Scan(&categoryID)

	if err == sql.ErrNoRows {
		// log.Printf("[collector] category not found: %s", categorySlug)
		slog.Error(" no category not found:", "category_slug", categorySlug)

		http.Error(w, "category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		// log.Printf("[collector] db error looking up category: %v", err)
		slog.Error("db error looking up category", "err", err) // Old Logging
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Step 3 — get amenities for this category
	amenities := config.AmenitiesForCategory(categorySlug)
	if len(amenities) == 0 {
		log.Printf("[collector] no amenities configured for category: %s", categorySlug)
		http.Error(w, "no amenities configured for this category", http.StatusUnprocessableEntity)
		return
	}

	// Step 4 — fetch from Overpass
	places, err := overpass.FetchPlaces(cityID, categoryID, amenities, bbox)
	if err != nil {
		// log.Printf("[collector] fetch failed: %v", err)
		http.Error(w, "collection failed", http.StatusInternalServerError)
		return
	}

	// Step 5 — insert into DB
	inserted, err := db.InsertPlaces(database, places)
	if err != nil {
		// log.Printf("[collector] insert failed: %v", err)
		http.Error(w, "database write failed", http.StatusInternalServerError)
		return
	}

	// Old Logging
	// log.Printf("[collector] complete: city=%s category=%s fetched=%d inserted=%d", citySlug, categorySlug, len(places), inserted)
	slog.Info("collection completed",
		"city_name", citySlug, // note: city_name not city
		"cat", categorySlug, // note: cat not category
		"fetched", len(places),
		"inserted", inserted,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"city":     citySlug,
		"category": categorySlug,
		"fetched":  len(places),
		"inserted": inserted,
	})
}
