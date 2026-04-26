package db

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/city-explorer-collector/internal/models"
)

// InsertPlace writes a single place to the database.
// It uses ON CONFLICT DO NOTHING to skip duplicates
// based on source + source_id.
// This makes the insert safe to run multiple times.
func InsertPlace(database *sql.DB, place models.Place) error {
	query := `
		INSERT INTO places (
			city_id, category_id, name, subcategory,
			address, latitude, longitude,
			website, phone, source, source_id
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10, $11
		)
		ON CONFLICT (source, source_id) DO NOTHING
	`

	_, err := database.Exec(query,
		place.CityID,
		place.CategoryID,
		place.Name,
		place.Subcategory,
		place.Address,
		place.Latitude,
		place.Longitude,
		place.Website,
		place.Phone,
		place.Source,
		place.SourceID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert place %s: %w", place.Name, err)
	}

	return nil
}

// InsertPlaces writes a slice of places to the database.
// It logs how many were inserted vs skipped.
func InsertPlaces(database *sql.DB, places []models.Place) (int, error) {
	inserted := 0

	for _, place := range places {
		err := InsertPlace(database, place)
		if err != nil {
			return inserted, err
		}
		inserted++
	}

	return inserted, nil
}
