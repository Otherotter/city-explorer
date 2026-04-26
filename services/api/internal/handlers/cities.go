package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/yourusername/city-explorer-api/internal/collector"
	apidb "github.com/yourusername/city-explorer-api/internal/db"
)

type PlaceResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Subcategory  string  `json:"subcategory"`
	Address      string  `json:"address"`
	Neighborhood string  `json:"neighborhood"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Website      string  `json:"website"`
	Phone        string  `json:"phone"`
	Source       string  `json:"source"`
}

type CitiesHandler struct {
	DB              *sql.DB
	CollectorClient *collector.Client
}

func (h *CitiesHandler) FoodHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("city-explorer-api")
	ctx := r.Context()

	cityName := chi.URLParam(r, "name")
	cityName = strings.ToLower(strings.TrimSpace(cityName))

	// Phase 1 field name: city (API uses "city", collector uses "city_name")
	// This inconsistency is intentional for Phase 1
	slog.Info("food request received",
		"city", cityName,
	)

	// Span 1 — city lookup
	ctx, citySpan := tracer.Start(ctx, "db.cities.lookup")
	citySpan.SetAttributes(attribute.String("city.name", cityName))

	var cityID int
	err := h.DB.QueryRowContext(ctx,
		"SELECT id FROM cities WHERE LOWER(name) = $1",
		cityName,
	).Scan(&cityID)

	citySpan.End()

	if err == sql.ErrNoRows {
		slog.Warn("city not found", "city", cityName)
		http.Error(w, "city not found", http.StatusNotFound)
		return
	}
	if err != nil {
		// log.Printf("[api] error looking up city: %v", err)
		slog.Error("db error looking up city",
			"city", cityName,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Span 2 — check if data is fresh
	ctx, checkSpan := tracer.Start(ctx, "db.places.check_status")
	checkSpan.SetAttributes(
		attribute.Int("city.id", cityID),
		attribute.String("category", "food"),
	)

	status, err := apidb.CheckDataStatus(h.DB, cityID, "food", 30)
	checkSpan.End()

	if err != nil {
		// log.Printf("[api] error checking data status: %v", err)
		slog.Error("error checking data status",
			"city", cityName,
			"city_id", cityID,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("data status checked",
		"city", cityName,
		"is_empty", status.IsEmpty,
		"is_stale", status.IsStale,
		"count", status.Count,
	)

	// Span 3 — trigger collection if needed
	if status.IsEmpty || status.IsStale {
		// log.Printf("[api] data missing or stale for %s/food, triggering collection", cityName)
		slog.Info("triggering collection",
			"city", cityName,
			"category", "food",
			"reason_empty", status.IsEmpty,
			"reason_stale", status.IsStale,
		)

		ctx, collectSpan := tracer.Start(ctx, "collector.trigger")
		collectSpan.SetAttributes(
			attribute.String("city", cityName),
			attribute.String("category", "food"),
			attribute.Bool("is_empty", status.IsEmpty),
			attribute.Bool("is_stale", status.IsStale),
		)

		err := h.CollectorClient.Collect(ctx, cityName, "food")
		collectSpan.End()

		if err != nil {
			// Collection failed but we can still try to return
			// whatever data we have in the DB
			// log.Printf("[api] collection failed: %v", err)

			slog.Warn("collection failed, returning existing data",
				"city", cityName,
				"category", "food",
				"error", err,
			)
		}
	}

	// Span 4 — fetch places from DB
	ctx, placesSpan := tracer.Start(ctx, "db.places.fetch_by_category")
	placesSpan.SetAttributes(
		attribute.Int("city.id", cityID),
		attribute.String("category", "food"),
	)

	rows, err := h.DB.QueryContext(ctx, `
		SELECT
			p.id,
			p.name,
			COALESCE(p.subcategory, '')  as subcategory,
			COALESCE(p.address, '')      as address,
			COALESCE(p.neighborhood, '') as neighborhood,
			p.latitude,
			p.longitude,
			COALESCE(p.website, '')      as website,
			COALESCE(p.phone, '')        as phone,
			p.source
		FROM places p
		JOIN categories c ON p.category_id = c.id
		WHERE p.city_id = $1
		AND c.slug = 'food'
		ORDER BY p.name ASC
		LIMIT 100
	`, cityID)

	placesSpan.End()

	if err != nil {
		slog.Error("db error fetching food places",
			"city", cityName,
			"city_id", cityID,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	places := []PlaceResponse{}

	for rows.Next() {
		var p PlaceResponse
		err := rows.Scan(
			&p.ID, &p.Name, &p.Subcategory,
			&p.Address, &p.Neighborhood,
			&p.Latitude, &p.Longitude,
			&p.Website, &p.Phone, &p.Source,
		)
		if err != nil {
			// log.Printf("[api] error scanning row: %v", err)
			slog.Warn("error scanning place row", "error", err)
			continue
		}
		places = append(places, p)
	}

	if err := rows.Err(); err != nil {
		slog.Error("row iteration error",
			"city", cityName,
			"error", err,
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// log.Printf("[api] returning %d food places for %s", len(places), cityName)
	slog.Info("returning food places",
		"city", cityName,
		"count", len(places),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"city":   cityName,
		"count":  len(places),
		"places": places,
	})
}
