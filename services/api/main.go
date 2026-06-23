package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Otherotter/city-explorer/services/api/internal/collector"
	"github.com/Otherotter/city-explorer/shared/observability"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Otherotter/city-explorer/services/api/internal/db"
	"github.com/Otherotter/city-explorer/services/api/internal/handlers"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

var database *sql.DB

func main() {
	ctx := context.Background()

	// REPLACE with this
	slog.SetDefault(observability.NewLogger("api"))
	slog.Info("api service starting...")

	shutdown, traceerr := observability.InitTracer(ctx, "city-explorer-api")
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

	// Connect to database
	var err error
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

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	// No per-route otelhttp needed anymore
	// The global handler below covers everything
	r.Get("/health", healthHandler)
	r.Handle("/metrics", promhttp.Handler())
	// r.Get("/cities/{name}/food", citiesHandler.FoodHandler)
	// City routes — traced individually
	// Apply otelhttp per route group
	// This correctly propagates context into handlers

	r.Route("/cities/{name}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "cities",
				otelhttp.WithMessageEvents(
					otelhttp.ReadEvents,
					otelhttp.WriteEvents,
				),
			)
		})
		r.Get("/food", citiesHandler.FoodHandler)
		// Add other categories here as you build them
		// r.Get("/nature", citiesHandler.NatureHandler)
		// r.Get("/study", citiesHandler.StudyHandler)
	})

	port := getEnv("PORT", "8080")
	slog.Info("server started", "port", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
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
