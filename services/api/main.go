package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Otherotter/city-explorer/services/api/internal/collector"
	"github.com/Otherotter/city-explorer/shared/observability"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Otherotter/city-explorer/services/api/internal/db"
	"github.com/Otherotter/city-explorer/services/api/internal/handlers"
	"github.com/Otherotter/city-explorer/services/api/internal/telemetry"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

var database *sql.DB

func main() {
	// REPLACE with this
	slog.SetDefault(observability.NewLogger("api"))
	slog.Info("api service starting...")

	ctx := context.Background()
	// Try local dev path first, fall back silently
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}
	slog.Info("environment loaded")

	// Initialize tracer // Must happen before any handlers are set up
	shutdown, err := telemetry.InitTracer(ctx, "city-explorer-api")
	if err != nil {
		// log.Fatalf("[api] failed to init tracer: %v", err) // Old Logging
		slog.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer shutdown(ctx)

	slog.Info("tracer initialized") // New Logging
	// log.Println("[api] tracer initialized") // Old Logging

	// Connect to database
	database, err = db.Connect()
	if err != nil {
		// log.Fatalf("[api] could not connect to database: %v", err)
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	// log.Println("[api] database connected") // Old Logging
	slog.Info("database connected") // New Logging

	collectorClient := collector.NewClient() // Initialize collector client
	// Initialize handlers with DB and collector client
	citiesHandler := &handlers.CitiesHandler{
		DB:              database,
		CollectorClient: collectorClient,
	}

	r := chi.NewRouter()
	// Replace middleware.Logger with:
	r.Use(chimiddleware.Recoverer)
	r.Use(observability.RequestLogger)

	// System routes — no tracing needed
	r.Get("/health", healthHandler)
	r.Handle("/metrics", promhttp.Handler())

	// City routes — wrapped with otelhttp
	// This automatically creates a trace span
	// for every request to these routes
	r.With(otelhttp.NewMiddleware("cities")).
		Get("/cities/{name}/food", citiesHandler.FoodHandler)

	port := getEnv("PORT", "8080")
	// log.Printf("[api] server starting on port %s", port) //Old Logging
	slog.Info("server started", "port", port) // New Logging

	if err := http.ListenAndServe(":"+port, r); err != nil {
		// log.Fatalf("[api] server failed: %v", err)
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err := database.Ping(); err != nil {
		response.Status = "degraded"
		response.Database = "unreachable"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		response.Status = "ok"
		response.Database = "connected"
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
